# Microservices based sample application

With this sample application, you can test tracing between multiple services. This is a companion repository for the article published [here](https://medium.com/@carlosedp/instrumenting-go-for-tracing-c5bdabe1fc81?source=friends_link&sk=904ae1e0e00555c26a51ac3e2a138cea).

Run Jaeger Tracing all-in-one in your own machine by using Docker:

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
  jaegertracing/all-in-one:latest
```

Open Jaeger GUI by using the URL [http://localhost:16686](http://localhost:16686).

## Building

To build all four binaries, just run `make`.

## Running

Start four different shells, run each server binary (`publisher`, `formatter`, `echo-app`) on it's own shell to see logs.

Run the `client` app with a name as argument like `./client Joe`.

Check traces on Jaeger GUI.
