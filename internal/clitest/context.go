package clitest

import (
	"errors"
	"testing"

	"github.com/hansmi/papercli/internal/cli"
	plclient "github.com/hansmi/paperhooks/pkg/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type Context struct {
	logger *zap.Logger
}

var _ cli.Context = (*Context)(nil)

func NewContext(t *testing.T) *Context {
	t.Helper()

	c := &Context{
		logger: zaptest.NewLogger(t),
	}

	return c
}

func (c *Context) Logger() *zap.Logger {
	return c.logger
}

func (c *Context) Client() (*plclient.Client, error) {
	return nil, errors.New("client not implemented")
}
