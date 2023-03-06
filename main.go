package main

import (
	"demo/fun"
	"fmt"
	"os"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/pkg/errors"
)

func main() {
	lg, _ := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://c3a49d7bfa894fba812faf3f4c0f1c89@o4504353388494848.ingest.sentry.io/4504789533065216",
		Debug:            true,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
	})
	if err != nil {
		lg.Sugar().Errorf("sentry.Init: %s", err)
	}

	sentryCfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, // when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
	}

	if sentry.CurrentHub().Client() == nil {
		return
	}
	core, _ := zapsentry.NewCore(sentryCfg, zapsentry.NewSentryClientFromClient(sentry.CurrentHub().Client()))

	lg = zapsentry.AttachCoreToLogger(core, lg)
	// lg.Error("sentry test")

	// sentry.CaptureException(errors.New("test error"))

	var result *multierror.Error

	result = multierror.Append(result, errors.New("new error 1111"))
	if err := fun.OpenFile1(); err != nil {
		result = multierror.Append(result, err)
		// lg.Error("OpenFile", zap.Error(err))
	}

	result = multierror.Append(result, errors.New("new error 2222"))
	if err := fun.OpenFile2(); err != nil {
		result = multierror.Append(result, err)
		// lg.Error("OpenFile", zap.Error(err))
	}
	result = multierror.Append(result, parseArgs(os.Args[1:]))

	info := make(map[string]interface{})
	for _, err := range result.WrappedErrors() {
		info[err.Error()] = fmt.Sprintf("%+v\n", err)
	}
	lg.Sugar().Errorw("multierror test", "error", info)
	time.Sleep(2 * time.Second)
}

func parseArgs(args []string) error {
	if len(args) < 3 {
		return errors.Errorf("not enough arguments, expected at least 3, got %d", len(args))
	}
	return nil
}
