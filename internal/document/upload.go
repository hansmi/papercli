package document

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/papercli/internal/cli"
	plclient "github.com/hansmi/paperhooks/pkg/client"
	"go.uber.org/zap"
)

type uploadClient interface {
	UploadDocument(context.Context, io.Reader, plclient.DocumentUploadOptions) (*plclient.DocumentUpload, *plclient.Response, error)
	WaitForTask(context.Context, string, plclient.WaitForTaskOptions) (*plclient.Task, error)
}

type uploadHandler struct {
	path         string
	filename     string
	wait         bool
	waitDuration time.Duration

	client uploadClient
}

func (h *uploadHandler) Setup(cmd *kingpin.CmdClause) {
	cmd.Help("Upload a document for consumption.")

	cmd.Arg("path", "Path to input file.").
		Required().
		StringVar(&h.path)

	cmd.Flag("filename", "Override provided filename.").
		StringVar(&h.filename)
	cmd.Flag("wait", "Wait for document to be consumed.").
		BoolVar(&h.wait)
	cmd.Flag("wait_duration", "Maximum amount of time to wait.").
		Default("1h").
		DurationVar(&h.waitDuration)
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

	result, _, err := client.UploadDocument(ctx, fh, opts)
	if err != nil {
		return fmt.Errorf("uploading document failed: %w", err)
	}

	logger.Info("Document uploaded successfully",
		zap.String("task_id", result.TaskID))

	if task, err := client.WaitForTask(ctx, result.TaskID, plclient.WaitForTaskOptions{
		MaxElapsedTime: h.waitDuration,
	}); err != nil {
		return fmt.Errorf("waiting for task: %w", err)
	} else {
		logger.Info("Task finished", zap.Any("task", task))
	}

	return nil
}
