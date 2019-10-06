package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) (opentracing.Tracer, io.Closer) {
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

// TraceFunction wraps funtion with opentracing span adding tags for the function name and caller details
func TraceFunction(ctx context.Context, fn interface{}, params ...interface{}) (result []reflect.Value) {
	// Get function name
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Create child span
	parentSpan := opentracing.SpanFromContext(ctx)
	sp := opentracing.StartSpan(
		"Function - "+name,
		opentracing.ChildOf(parentSpan.Context()))
	defer sp.Finish()

	sp.SetTag("function", name)

	// Get caller function name, file and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	callerDetails := fmt.Sprintf("%s - %s#%d", frame.Function, frame.File, frame.Line)
	sp.SetTag("caller", callerDetails)

	// Check params and call function
	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(params) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(params))
	for k, in := range params {
		inputs[k] = reflect.ValueOf(in)
	}
	return f.Call(inputs)
}

// CreateChildSpan creates a new opentracing span adding tags for the span name and caller details. Returns a Span.
// User must call `defer sp.Finish()`
func CreateChildSpan(ctx context.Context, name string) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(ctx)
	sp := opentracing.StartSpan(
		name,
		opentracing.ChildOf(parentSpan.Context()))
	sp.SetTag("name", name)

	// Get caller function name, file and line
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	callerDetails := fmt.Sprintf("%s - %s#%d", frame.Function, frame.File, frame.Line)
	sp.SetTag("caller", callerDetails)

	return sp
}

// NewTracedRequest generates a new traced HTTP request with opentracing headers injected into it
func NewTracedRequest(method string, url string, body io.Reader, span opentracing.Span) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, method)
	span.Tracer().Inject(span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	return req, err
}
