package document

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/papercli/internal/cli"
	plclient "github.com/hansmi/paperhooks/pkg/client"
)

type fakeUploadClient struct {
	listTagsErr    error
	waitForTaskErr error
}

func (c *fakeUploadClient) ListTags(ctx context.Context, opts plclient.ListTagsOptions) ([]plclient.Tag, *plclient.Response, error) {
	return nil, nil, c.listTagsErr
}

func (c *fakeUploadClient) UploadDocument(context.Context, io.Reader, plclient.DocumentUploadOptions) (*plclient.DocumentUpload, *plclient.Response, error) {
	return &plclient.DocumentUpload{
		TaskID: "upload-document-result",
	}, nil, nil
}

func (c *fakeUploadClient) WaitForTask(context.Context, string, plclient.WaitForTaskOptions) (*plclient.Task, error) {
	return &plclient.Task{
		TaskID: "wait-for-task-result",
	}, c.waitForTaskErr
}

func TestUploadHandler(t *testing.T) {
	errTest := errors.New("test error")
	errDuplicate := &plclient.TaskError{
		Status:  plclient.TaskFailure,
		Message: "Not consuming xyz.pdf: It is a duplicate of Name of another document (#1234)",
	}

	for _, tc := range []struct {
		name    string
		h       uploadHandler
		wantErr error
	}{
		{
			name: "missing file",
			h: uploadHandler{
				path: filepath.Join(t.TempDir(), "missing"),
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "success",
			h: uploadHandler{
				path: os.DevNull,
			},
		},
		{
			name: "tag not found",
			h: uploadHandler{
				path:     os.DevNull,
				tagNames: []string{"a"},
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "tag error",
			h: uploadHandler{
				path:     os.DevNull,
				tagNames: []string{"a"},
				client: &fakeUploadClient{
					listTagsErr: errTest,
				},
			},
			wantErr: errTest,
		},
		{
			name: "wait error",
			h: uploadHandler{
				path: os.DevNull,
				client: &fakeUploadClient{
					waitForTaskErr: errTest,
				},
			},
			wantErr: errTest,
		},
		{
			name: "duplicate document",
			h: uploadHandler{
				path: os.DevNull,
				client: &fakeUploadClient{
					waitForTaskErr: errDuplicate,
				},
			},
			wantErr: errDuplicate,
		},
		{
			name: "duplicate document suppressed",
			h: uploadHandler{
				path:            os.DevNull,
				ignoreDuplicate: true,
				client: &fakeUploadClient{
					waitForTaskErr: errDuplicate,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			t.Cleanup(cancel)

			if tc.h.client == nil {
				tc.h.client = &fakeUploadClient{}
			}

			err := tc.h.Run(ctx, cli.NewContextForTest(t))

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Run() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
