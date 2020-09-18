package tracer

import (
	"fmt"
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func InitSpan(name string, tracer opentracing.Tracer, parentSpan opentracing.Span) opentracing.Span {

	env := os.Getenv("ENVIRONMENT")
	user := os.Getenv("USERNAME")

	var span opentracing.Span

	if parentSpan == nil {
		span = tracer.StartSpan(name)
	} else {
		span = tracer.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
	}

	span.SetTag("env", env)
	span.SetTag("user", user)

	return span
}

// func SetTagsForSpan(method string, url string, statuscode string, span opentracing.Span) opentracing.Span {

// 	if strings.HasPrefix(statuscode, "5") {
// 		span.SetTag("statuscode", "5XX")
// 	} else if strings.HasPrefix(statuscode, "4") {
// 		span.SetTag("statuscode", "4XX")
// 	} else if strings.HasPrefix(statuscode, "3") {
// 		span.SetTag("statuscode", "3XX")
// 	} else {
// 		span.SetTag("statuscode", statuscode)
// 	}

// 	span.SetTag("url", url)
// 	span.SetTag("method", method)

// 	return span
// }
