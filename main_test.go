package main

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestBuildLoggerConfig(t *testing.T) {
	config := buildLoggerConfig()

	if _, err := config.Build(); err != nil {
		t.Errorf("Building logger failed: %v", err)
	}
}

func TestRegisterCommands(t *testing.T) {
	app := kingpin.New("", "")

	reg := registerCommands(app)

	for _, name := range []string{
		"document upload",
	} {
		if got := reg.Lookup(name); got == nil {
			t.Errorf("Missing command %q", name)
		}
	}
}

func TestProgramRun(t *testing.T) {
	for _, tc := range []struct {
		name    string
		args    []string
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: kingpin.ErrCommandNotSpecified,
		},
		{
			name:    "invalid command",
			args:    []string{"#never valid command"},
			wantErr: cmpopts.AnyError,
		},
		{
			name:    "invalid flag",
			args:    []string{"-\t"},
			wantErr: cmpopts.AnyError,
		},
		{
			name:    "incomplete flags",
			args:    []string{"document", "upload", os.DevNull},
			wantErr: cmpopts.AnyError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app := kingpin.New("", "")
			app.Terminate(nil)
			app.UsageWriter(io.Discard)

			p := newProgram(zap.NewProductionConfig())
			p.logger = zaptest.NewLogger(t)

			err := p.run(context.Background(), app, tc.args)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("run() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
