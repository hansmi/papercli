package document

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/papercli/internal/cli"
	plclient "github.com/hansmi/paperhooks/pkg/client"
	"go.uber.org/zap"
)

var taskResultNotConsumingDuplicateRe = regexp.MustCompile(`(?i)\bnot consuming\b.*: It is a duplicate of\b.*\(#\d+\)\s*(?:\.\s*)?$`)

type uploadClient interface {
	ListTags(context.Context, plclient.ListTagsOptions) ([]plclient.Tag, *plclient.Response, error)
	UploadDocument(context.Context, io.Reader, plclient.DocumentUploadOptions) (*plclient.DocumentUpload, *plclient.Response, error)
	WaitForTask(context.Context, string, plclient.WaitForTaskOptions) (*plclient.Task, error)
}

type uploadHandler struct {
	client uploadClient

	path string

	filename string
	tagNames []string

	wait         bool
	waitDuration time.Duration

	ignoreDuplicate bool
}

func newUploadHandler() *uploadHandler {
	return &uploadHandler{
		wait:            true,
		waitDuration:    time.Hour,
		ignoreDuplicate: true,
	}
}

func (h *uploadHandler) Setup(cmd *kingpin.CmdClause) {
	cmd.Help("Upload a document for consumption.")

	cmd.Arg("path", "Path to input file.").
		Required().
		StringVar(&h.path)

	cmd.Flag("filename", "Override provided filename.").
		StringVar(&h.filename)
	cmd.Flag("tag", "Apply a pre-existing tag to the document.").
		StringsVar(&h.tagNames)

	cmd.Flag("wait", "Wait for document to be consumed.").
		BoolVar(&h.wait)
	cmd.Flag("wait_duration", "Maximum amount of time to wait.").
		DurationVar(&h.waitDuration)

	cmd.Flag("ignore_duplicate", "Suppress error status for duplicated documents.").
		BoolVar(&h.ignoreDuplicate)
}

func (h *uploadHandler) resolveTags(ctx context.Context, client uploadClient) ([]int64, error) {
	var result []int64

	for _, name := range h.tagNames {
		opts := plclient.ListTagsOptions{}
		opts.Name.EqualsIgnoringCase = &name

		if items, _, err := client.ListTags(ctx, opts); err != nil {
			return nil, fmt.Errorf("tag %q: %w", name, err)
		} else if len(items) == 0 {
			return nil, fmt.Errorf("tag %q not found", name)
		} else if len(items) != 1 {
			return nil, fmt.Errorf("tag %q: received %d items, expected exactly one", name, len(items))
		} else {
			result = append(result, items[0].ID)
		}
	}

	return result, nil
}

func (h *uploadHandler) Run(ctx context.Context, cctx cli.Context) error {
	logger := cctx.Logger()

	fh, err := os.Open(h.path)
	if err != nil {
		return err
	}

	defer fh.Close()

	client := h.client

	if client == nil {
		if client, err = cctx.Client(); err != nil {
			return err
		}
	}

	opts := plclient.DocumentUploadOptions{
		Filename: h.filename,
	}

	if opts.Filename == "" {
		opts.Filename = filepath.Base(h.path)
	}

	if opts.Tags, err = h.resolveTags(ctx, client); err != nil {
		return fmt.Errorf("resolving tags: %w", err)
	}

	result, _, err := client.UploadDocument(ctx, fh, opts)
	if err != nil {
		return fmt.Errorf("uploading document failed: %w", err)
	}

	logger.Info("Document uploaded successfully",
		zap.String("task_id", result.TaskID))

	if h.wait {
		task, err := client.WaitForTask(ctx, result.TaskID, plclient.WaitForTaskOptions{
			MaxElapsedTime: h.waitDuration,
		})
		if err != nil {
			var taskErr *plclient.TaskError

			if h.ignoreDuplicate && errors.As(err, &taskErr) && taskErr.Status == plclient.TaskFailure && taskResultNotConsumingDuplicateRe.MatchString(taskErr.Message) {
				logger.Error("Document is a duplicate", zap.Error(taskErr))
				return nil
			}

			return fmt.Errorf("waiting for document consumption: %w", err)
		}

		logger.Info("Document consumed successfully", zap.Any("task", task))
	}

	return nil
}
