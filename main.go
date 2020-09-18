package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("init started")

}

func main() {

	finish := make(chan bool)

	tracer, closer := InitJaeger("sample-application")
	defer closer.Close()

	go func() {
		for count := 0; ; count++ {
			log.Println("Starting iteration ::: ", count)
			go sampleTest(tracer)
			time.Sleep(10 * time.Second)

		}
	}()

	<-finish
}

func sampleTest(tracer opentracing.Tracer) {

	span := InitSpan("parent span", tracer, nil)
	spanLogin := InitSpan("login", tracer, span)
	spanLogin.LogEventWithPayload("Input-credentilals", "sample@gmail.com samplepassword")
	spanLogin.LogEventWithPayload("Output-response", "OutputResponse")
	spanLogin.Finish()
}

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
