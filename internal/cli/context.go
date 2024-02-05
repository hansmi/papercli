package cli

import (
	"errors"
	"fmt"
	"testing"

	plclient "github.com/hansmi/paperhooks/pkg/client"
	"github.com/hansmi/paperhooks/pkg/kpflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type Context interface {
	Logger() *zap.Logger
	Client() (*plclient.Client, error)
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

type contextTestImpl struct {
	logger      *zap.Logger
	clientFlags plclient.Flags
}

func NewContextForTest(t *testing.T) Context {
	t.Helper()

	c := &contextTestImpl{
		logger: zaptest.NewLogger(t),
	}

	return c
}

func (c *contextTestImpl) Logger() *zap.Logger {
	return c.logger
}

func (c *contextTestImpl) Client() (*plclient.Client, error) {
	return nil, errors.New("client not implemented")
}
