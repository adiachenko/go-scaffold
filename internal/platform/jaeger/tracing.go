package jaeger

import (
	"adiachenko/go-scaffold/internal/platform/conf"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/noop"
	"github.com/Vinelab/tracing-go/drivers/zipkin"
	"github.com/sirupsen/logrus"
)

var (
	Trace tracing.Tracer
)

func init() {
	var err error

	switch conf.TracingDriver {
	case "zipkin":
		Trace, err = zipkin.NewTracer(zipkin.TracerOptions{
			ServiceName: conf.TracingServiceName,
			Host:        conf.ZipkinHost,
			Port:        conf.ZipkinPort,
		})
	case "noop":
		Trace = noop.NewTracer()
	default:
		Trace = noop.NewTracer()
	}

	if err != nil {
		logrus.WithError(err).Fatal(err.Error())
	}
}
