package main

import (
	"fmt"
	"log"
	"net/http"

	xhttp "microservices/lib/http"
	"microservices/lib/tracing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	tracer, closer := tracing.Init("publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
		url := "http://localhost:8080/test/" + helloStr
		req, err := tracing.NewTracedRequest("GET", url, nil, span)

		resp, err := xhttp.Do(req)
		if err != nil {
			log.Println("Could not contact echo: " + err.Error())
			span.LogEvent("Could not contact echo: " + err.Error())
			span.SetTag("error", true)
			return
		}
		fmt.Println(string(resp))
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
