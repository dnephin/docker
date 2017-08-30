package image

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestImageRouterDeleteImages(t *testing.T) {
	var testcases = []struct {
		doc          string
		vars         map[string]string
		expectedErr  string
		expectedResp []types.ImageDeleteResponseItem
	}{
		{
			doc: "no name param",
		},
		{
			doc: "error from image delete",
		},
		{
			doc: "successful delete",
		},
	}

	for _, testcase := range testcases {
		backend := &fakeBackend{}
		router := NewRouter(backend, nil).(*imageRouter)
		err := router.deleteImages(_, writer, request, testcase.vars)
	}
}
