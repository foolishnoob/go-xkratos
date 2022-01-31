package tracer

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(NewTracer)
	var err = container.Invoke(func(t Tracer) {
		otel.SetTracerProvider(t.GetInstance())
	})
	xdebug.IfPanic(err)
}

type Tracer interface {
	GetInstance() *traceSdk.TracerProvider
}

type tracer struct {
	traceProvider *traceSdk.TracerProvider
}

func (t *tracer) GetInstance() *traceSdk.TracerProvider {
	return t.traceProvider
}

func NewTracer(conf *config.BootConfig) Tracer {
	if nil != conf.GetTrace() && "" != conf.GetTrace().GetEndpoint() {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.GetTrace().GetEndpoint())))
		xdebug.IfPanic(err)
		return &tracer{
			traceProvider: traceSdk.NewTracerProvider(
				traceSdk.WithBatcher(exp),
				traceSdk.WithResource(resource.NewSchemaless(
					semConv.ServiceNameKey.String(conf.GetService().GetName()),
				)),
			),
		}
	}
	return nil
}
