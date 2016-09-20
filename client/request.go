package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// ServerResponse is a wrapper for http API responses.
type ServerResponse struct {
	Body       io.ReadCloser
	Header     http.Header
	StatusCode int
}

// Request is a data object for an HTTP request to the Docker API
type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Body    io.Reader
	Headers map[string][]string
}

// head sends an http request to the docker API using the method HEAD.
func (cli *Client) head(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.Do(ctx, Request{"HEAD", path, query, nil, headers})
}

// getWithContext sends an http request to the docker API using the method GET with a specific go context.
func (cli *Client) get(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.Do(ctx, Request{"GET", path, query, nil, headers})
}

// postWithContext sends an http request to the docker API using the method POST with a specific go context.
func (cli *Client) post(ctx context.Context, path string, query url.Values, obj interface{}, headers map[string][]string) (ServerResponse, error) {
	return cli.DoWithBody(ctx, Request{"POST", path, query, nil, headers}, obj)
}

func (cli *Client) postRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (ServerResponse, error) {
	return cli.Do(ctx, Request{"POST", path, query, body, headers})
}

// put sends an http request to the docker API using the method PUT.
func (cli *Client) put(ctx context.Context, path string, query url.Values, obj interface{}, headers map[string][]string) (ServerResponse, error) {
	return cli.DoWithBody(ctx, Request{"PUT", path, query, nil, headers}, obj)
}

// put sends an http request to the docker API using the method PUT.
func (cli *Client) putRaw(ctx context.Context, path string, query url.Values, body io.Reader, headers map[string][]string) (ServerResponse, error) {
	return cli.Do(ctx, Request{"PUT", path, query, body, headers})
}

// delete sends an http request to the docker API using the method DELETE.
func (cli *Client) delete(ctx context.Context, path string, query url.Values, headers map[string][]string) (ServerResponse, error) {
	return cli.Do(ctx, Request{"DELETE", path, query, nil, headers})
}

// DoWithBody encodes obj to JSON then performs a Request
func (cli *Client) DoWithBody(ctx context.Context, req Request, obj interface{}) (ServerResponse, error) {
	if obj != nil {
		var err error
		req.Body, err = encodeData(obj)
		if err != nil {
			return ServerResponse{}, err
		}
		if req.Headers == nil {
			req.Headers = make(map[string][]string)
		}
		req.Headers["Content-Type"] = []string{"application/json"}
	}

	return cli.Do(ctx, req)
}

// Do performs a request to the API
func (cli *Client) Do(ctx context.Context, apireq Request) (ServerResponse, error) {
	serverResp := ServerResponse{
		Body:       nil,
		StatusCode: -1,
	}

	expectedPayload := (apireq.Method == "POST" || apireq.Method == "PUT")
	if expectedPayload && apireq.Body == nil {
		apireq.Body = bytes.NewReader([]byte{})
	}

	req, err := cli.newRequest(apireq)
	if err != nil {
		return serverResp, err
	}

	if cli.proto == "unix" || cli.proto == "npipe" {
		// For local communications, it doesn't matter what the host is. We just
		// need a valid and meaningful host name. (See #189)
		req.Host = "docker"
	}

	scheme, err := resolveScheme(cli.client.Transport)
	if err != nil {
		return serverResp, err
	}

	req.URL.Host = cli.addr
	req.URL.Scheme = scheme

	if expectedPayload && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/plain")
	}

	resp, err := ctxhttp.Do(ctx, cli.client, req)
	if err != nil {

		if scheme == "https" && strings.Contains(err.Error(), "malformed HTTP response") {
			return serverResp, fmt.Errorf("%v.\n* Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}

		if scheme == "https" && strings.Contains(err.Error(), "bad certificate") {
			return serverResp, fmt.Errorf("The server probably has client authentication (--tlsverify) enabled. Please check your TLS client certification settings: %v", err)
		}

		// Don't decorate context sentinel errors; users may be comparing to
		// them directly.
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			return serverResp, err
		}

		if err, ok := err.(net.Error); ok {
			if err.Timeout() {
				return serverResp, ErrorConnectionFailed(cli.host)
			}
			if !err.Temporary() {
				if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial unix") {
					return serverResp, ErrorConnectionFailed(cli.host)
				}
			}
		}

		return serverResp, errors.Wrap(err, "error during connect")
	}

	if resp != nil {
		serverResp.StatusCode = resp.StatusCode
	}

	if serverResp.StatusCode < 200 || serverResp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return serverResp, err
		}
		if len(body) == 0 {
			return serverResp, fmt.Errorf("Error: request returned %s for API route and version %s, check if the server supports the requested API version", http.StatusText(serverResp.StatusCode), req.URL)
		}

		var errorMessage string
		if (cli.version == "" || versions.GreaterThan(cli.version, "1.23")) &&
			resp.Header.Get("Content-Type") == "application/json" {
			var errorResponse types.ErrorResponse
			if err := json.Unmarshal(body, &errorResponse); err != nil {
				return serverResp, fmt.Errorf("Error reading JSON: %v", err)
			}
			errorMessage = errorResponse.Message
		} else {
			errorMessage = string(body)
		}

		return serverResp, fmt.Errorf("Error response from daemon: %s", strings.TrimSpace(errorMessage))
	}

	serverResp.Body = resp.Body
	serverResp.Header = resp.Header
	return serverResp, nil
}

func (cli *Client) newRequest(req Request) (*http.Request, error) {
	apiPath := cli.getAPIPath(req.Path, req.Query)
	httpreq, err := http.NewRequest(req.Method, apiPath, req.Body)
	if err != nil {
		return nil, err
	}

	// Add CLI Config's HTTP Headers BEFORE we set the Docker headers
	// then the user can't change OUR headers
	for k, v := range cli.customHTTPHeaders {
		httpreq.Header.Set(k, v)
	}

	if req.Headers != nil {
		for k, v := range req.Headers {
			httpreq.Header[k] = v
		}
	}

	return httpreq, nil
}

func encodeData(data interface{}) (*bytes.Buffer, error) {
	params := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(params).Encode(data); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func ensureReaderClosed(response ServerResponse) {
	if body := response.Body; body != nil {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, body, 512)
		response.Body.Close()
	}
}
