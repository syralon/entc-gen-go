package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"{{.module}}/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(hs *http.Server, gs *grpc.Server) *kratos.App {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

// Set global trace provider
func setTracerProvider(trace *conf.Trace) error {
	opts := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	}
	if trace != nil {
		// Create the Jaeger exporter
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(trace.Endpoint)))
		if err != nil {
			return err
		}
		// Always be sure to batch in production.
		opts = append(opts, tracesdk.WithBatcher(exp))
	}

	tp := tracesdk.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	flag.Parse()

	cfg := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	if err := setTracerProvider(bc.Trace); err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.Server, bc.Data)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
