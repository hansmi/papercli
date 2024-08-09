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

func newUploadHandlerForTest(patch func(*uploadHandler)) *uploadHandler {
	h := newUploadHandler()
	h.path = os.DevNull
	h.client = &fakeUploadClient{}

	if patch != nil {
		patch(h)
	}

	return h
}

func TestUploadHandler(t *testing.T) {
	errTest := errors.New("test error")
	errDuplicateLegacy := &plclient.TaskError{
		Status:  plclient.TaskFailure,
		Message: "Not consuming xyz.pdf: It is a duplicate of Name of another document (#1234)",
	}
	errDuplicate := &plclient.TaskError{
		Status:  plclient.TaskFailure,
		Message: "Not consuming xyz.pdf: It is a duplicate of Name of another document (#1234).",
	}
	errDuplicateInTrash := &plclient.TaskError{
		Status:  plclient.TaskFailure,
		Message: "Not consuming xyz.pdf: It is a duplicate of Name of another document (#1234). Note: existing document is in the trash.",
	}

	for _, tc := range []struct {
		name    string
		h       *uploadHandler
		wantErr error
	}{
		{
			name: "missing file",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.path = filepath.Join(t.TempDir(), "missing")
			}),
			wantErr: os.ErrNotExist,
		},
		{
			name: "success",
			h:    newUploadHandlerForTest(nil),
		},
		{
			name: "tag not found",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.tagNames = []string{"a"}
			}),
			wantErr: cmpopts.AnyError,
		},
		{
			name: "tag error",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.tagNames = []string{"a"}
				h.client = &fakeUploadClient{
					listTagsErr: errTest,
				}
			}),
			wantErr: errTest,
		},
		{
			name: "wait disabled",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.wait = false
				h.client = &fakeUploadClient{
					waitForTaskErr: errTest,
				}
			}),
		},
		{
			name: "wait error",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.client = &fakeUploadClient{
					waitForTaskErr: errTest,
				}
			}),
			wantErr: errTest,
		},
		{
			name: "duplicate document before 2.11.3",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.client = &fakeUploadClient{
					waitForTaskErr: errDuplicateLegacy,
				}
			}),
		},
		{
			name: "duplicate document after 2.11.3",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.client = &fakeUploadClient{
					waitForTaskErr: errDuplicate,
				}
			}),
		},
		{
			name: "duplicate document is in trash",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.client = &fakeUploadClient{
					waitForTaskErr: errDuplicateInTrash,
				}
			}),
			wantErr: errDuplicateInTrash,
		},
		{
			name: "duplicate document causes error",
			h: newUploadHandlerForTest(func(h *uploadHandler) {
				h.ignoreDuplicate = false
				h.client = &fakeUploadClient{
					waitForTaskErr: errDuplicate,
				}
			}),
			wantErr: errDuplicate,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			t.Cleanup(cancel)

			err := tc.h.Run(ctx, cli.NewContextForTest(t))

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Run() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
