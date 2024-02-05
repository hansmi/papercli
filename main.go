package main

import (
	"context"
	"fmt"
	stdlog "log"
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/papercli/internal/cli"
	"github.com/hansmi/papercli/internal/document"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// buildLoggerConfig returns a logger config writing to standard error.
func buildLoggerConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	cfg.EncoderConfig.EncodeName = func(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + loggerName + "]")
	}
	cfg.EncoderConfig.NewReflectedEncoder = nil
	cfg.Encoding = "console"

	return cfg
}

func registerCommands(app *kingpin.Application) *cli.CommandRegistry {
	reg := cli.NewCommandRegistry()

	document.RegisterCommands(reg, app)

	return reg
}

type program struct {
	logLevel zap.AtomicLevel
	logger   *zap.Logger
}

func newProgram(loggerConfig zap.Config) *program {
	p := &program{
		logLevel: loggerConfig.Level,
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		stdlog.Fatalf("Initializing logger failed: %v", err)
	}

	p.logger = logger

	return p
}

func (p *program) run(ctx context.Context, app *kingpin.Application, args []string) error {
	app.Help = "Command line client for Paperless-ngx."

	customLogLevel := zap.InfoLevel

	app.Flag("log_level", "Log level for stderr.").
		Default(customLogLevel.String()).
		SetValue(&customLogLevel)

	p.logLevel.SetLevel(customLogLevel)

	tbl := registerCommands(app)

	commandName, err := app.Parse(args)
	if err != nil {
		return err
	}

	cmd := tbl.Lookup(commandName)
	if cmd == nil {
		return fmt.Errorf("unknown subcommand %q", commandName)
	}

	return cmd.Run(ctx, cli.NewContext(p.logger, app))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p := newProgram(buildLoggerConfig())

	defer p.logger.Sync()
	defer zap.ReplaceGlobals(p.logger)()
	defer zap.RedirectStdLog(p.logger)()

	if err := p.run(context.Background(), kingpin.CommandLine, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}
