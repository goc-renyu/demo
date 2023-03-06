package main

import (
	"demo/fun"
	"errors"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	lg, _ := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://c3a49d7bfa894fba812faf3f4c0f1c89@o4504353388494848.ingest.sentry.io/4504789533065216",
		Debug:            true,
		TracesSampleRate: 1.0,
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

	lg.Sugar().Errorw("multierror", "error", result.WrappedErrors())
	time.Sleep(2 * time.Second)
}
