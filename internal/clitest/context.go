package clitest

import (
	"errors"
	"testing"

	"github.com/hansmi/papercli/internal/cli"
	plclient "github.com/hansmi/paperhooks/pkg/client"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type contextImpl struct {
	logger      *zap.Logger
	clientFlags plclient.Flags
}

func NewContext(t *testing.T) cli.Context {
	t.Helper()

	c := &contextImpl{
		logger: zaptest.NewLogger(t),
	}

	return c
}

func (c *contextImpl) Logger() *zap.Logger {
	return c.logger
}

func (c *contextImpl) Client() (*plclient.Client, error) {
	return nil, errors.New("client not implemented")
}
