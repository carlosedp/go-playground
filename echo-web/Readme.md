# Labstack Echo demo app using metrics and tracing

This project is a demo on how to use the Prometheus metrics and Jaeger tracing middleware
on a Labstack Echo Framework application.

To start Jaeger all-in-one tracer on a container in your own machine, run:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one
```
