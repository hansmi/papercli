package cli

import (
	"fmt"
	"io"
	"os"

	plclient "github.com/hansmi/paperhooks/pkg/client"
	"github.com/hansmi/paperhooks/pkg/kpflag"
	"go.uber.org/zap"
)

type Context interface {
	Logger() *zap.Logger
	Client() (*plclient.Client, error)
	Output() io.Writer
}

type contextImpl struct {
	logger      *zap.Logger
	clientFlags plclient.Flags
}

func NewContext(logger *zap.Logger, g kpflag.FlagGroup) Context {
	c := &contextImpl{
		logger: logger,
	}

	kpflag.RegisterClient(g, &c.clientFlags)

	return c
}

func (c *contextImpl) Logger() *zap.Logger {
	return c.logger
}

func (c *contextImpl) Client() (*plclient.Client, error) {
	client, err := c.clientFlags.Build()
	if err != nil {
		return nil, fmt.Errorf("building client: %w", err)
	}

	return client, nil
}

func (c *contextImpl) Output() io.Writer {
	return os.Stdout
}
