package dockerfile

import (
	"time"

	"github.com/docker/docker/builder/fscache"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/client/session"
	"github.com/docker/docker/client/session/filesync"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const sessionConnectTimeout = 7 * time.Second

// ClientSessionTransport is a transport for copying files from docker client
// to the daemon.
type ClientSessionTransport struct{}

// NewClientSessionTransport returns new ClientSessionTransport instance
func NewClientSessionTransport() *ClientSessionTransport {
	return &ClientSessionTransport{}
}

// Copy data from a remote to a destination directory.
func (cst *ClientSessionTransport) Copy(ctx context.Context, id fscache.RemoteIdentifier, dest string, cu filesync.CacheUpdater) error {
	csi, ok := id.(*ClientSessionIdentifier)
	if !ok {
		return errors.New("invalid identifier for client session")
	}

	return filesync.FSSync(ctx, csi.caller, filesync.FSSendRequestOpt{
		SrcPaths:     csi.srcPaths,
		DestDir:      dest,
		CacheUpdater: cu,
	})
}

// ClientSessionIdentifier is an identifier that can be used for requesting
// files from remote client
type ClientSessionIdentifier struct {
	srcPaths  []string
	caller    session.Caller
	sharedKey string
	uuid      string
}

// NewClientSessionIdentifier returns new ClientSessionIdentifier instance
func NewClientSessionIdentifier(ctx context.Context, sg SessionGetter, uuid string, sources []string) (*ClientSessionIdentifier, error) {
	csi := &ClientSessionIdentifier{
		uuid:     uuid,
		srcPaths: sources,
	}
	ctx, cancel := context.WithTimeout(ctx, sessionConnectTimeout)
	defer cancel()
	caller, err := sg.Get(ctx, uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get session for %s", uuid)
	}

	csi.caller = caller
	return csi, nil
}

// Transport returns transport identifier for remote identifier
func (csi *ClientSessionIdentifier) Transport() string {
	return remotecontext.ClientSessionRemote
}

// SharedKey returns shared key for remote identifier. Shared key is used
// for finding the base for a repeated transfer.
func (csi *ClientSessionIdentifier) SharedKey() string {
	return csi.caller.SharedKey()
}

// Key returns unique key for remote identifier. Requests with same key return
// same data.
func (csi *ClientSessionIdentifier) Key() string {
	return csi.uuid
}
