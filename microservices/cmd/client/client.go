package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	xhttp "microservices/lib/http"
	"microservices/lib/tracing"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

func main() {
	// if len(os.Args) != 2 {
	// 	log.Fatal("ERROR: Expecting one argument")
	// }
	// helloTo := os.Args[1]

	e := echo.New()
	c := jaegertracing.New(e, nil)
	defer c.Close()

	e.GET("/hello/:name", func(c echo.Context) error {
		var name = ""
		name = c.Param("name")
		out := fmt.Sprintf("Hello %s", name)
		helloStr := formatString(c.Request().Context(), name)
		printHello(c.Request().Context(), helloStr)
		return c.String(http.StatusOK, out)
	})
	e.Logger.Fatal(e.Start(":8080"))
	// tracer, closer := tracing.Init("hello-client")
	// defer closer.Close()
	// opentracing.SetGlobalTracer(tracer)

	// span := tracer.StartSpan("say-hello")
	// span.SetTag("hello-to", helloTo)
	// defer span.Finish()

	// ctx := opentracing.ContextWithSpan(context.Background(), span)

	// helloStr := formatString(ctx, helloTo)
	// printHello(ctx, helloStr)
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://formatter:8080/format?" + v.Encode()

	// Here we create a NewReques and manually inject the tracing headers to it
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := xhttp.Do(req)
	if err != nil {
		panic(err.Error())
	}

	helloStr := string(resp)

	span.LogFields(
		otlog.String("event", "string-format"),
		otlog.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
	defer span.Finish()

	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://publisher:8080/publish?" + v.Encode()

	// Here we use convenience function to create the request with tracing headers
	req, err := tracing.NewTracedRequest("GET", url, nil, span)
	if err != nil {
		panic(err.Error())
	}

	if _, err := xhttp.Do(req); err != nil {
		panic(err.Error())
	}
}
