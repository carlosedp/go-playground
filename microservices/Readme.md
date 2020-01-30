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
  -p 14267:14267 \
  -p 9411:9411 \
  jaegertracing/all-in-one
```

Open Jaeger GUI by using the URL [http://localhost:16686](http://localhost:16686).

## Building

To build all four binaries, just run `make`.

## Running

Start four different shells, run each server binary (`publisher`, `formatter`, `echo-app`) on it's own shell to see logs.

Run the `client` app with a name as argument like `./client Joe`.

Check traces on Jaeger GUI.

The tracer can be initialized with values coming from environment variables. None of the env vars are required
and all of them can be overriden via direct setting of the property on the configuration object.

Property| Description
--- | ---
JAEGER_SERVICE_NAME | The service name
JAEGER_AGENT_HOST | The hostname for communicating with agent via UDP
JAEGER_AGENT_PORT | The port for communicating with agent via UDP
JAEGER_ENDPOINT | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces
JAEGER_USER | Username to send as part of "Basic" authentication to the collector endpoint
JAEGER_PASSWORD | Password to send as part of "Basic" authentication to the collector endpoint
JAEGER_REPORTER_LOG_SPANS | Whether the reporter should also log the spans
JAEGER_REPORTER_MAX_QUEUE_SIZE | The reporter's maximum queue size
JAEGER_REPORTER_FLUSH_INTERVAL | The reporter's flush interval, with units, e.g. "500ms" or "2s" ([valid units][timeunits])
JAEGER_SAMPLER_TYPE | The sampler type
JAEGER_SAMPLER_PARAM | The sampler parameter (number)
JAEGER_SAMPLER_MANAGER_HOST_PORT | The HTTP endpoint when using the remote sampler, i.e. http://jaeger-agent:5778/sampling
JAEGER_SAMPLER_MAX_OPERATIONS | The maximum number of operations that the sampler will keep track of
JAEGER_SAMPLER_REFRESH_INTERVAL | How often the remotely controlled sampler will poll jaeger-agent for the appropriate sampling strategy, with units, e.g. "1m" or "30s" ([valid units][timeunits])
JAEGER_TAGS | A comma separated list of `name = value` tracer level tags, which get added to all reported spans. The value can also refer to an environment variable using the format `${envVarName:default}`, where the `:default` is optional, and identifies a value to be used if the environment variable cannot be found
JAEGER_DISABLED | Whether the tracer is disabled or not. If true, the default `opentracing.NoopTracer` is used.
JAEGER_RPC_METRICS | Whether to store RPC metrics

By default, the client sends traces via UDP to the agent at `localhost:6831`. Use `JAEGER_AGENT_HOST` and
`JAEGER_AGENT_PORT` to send UDP traces to a different `host:port`. If `JAEGER_ENDPOINT` is set, the client sends traces
to the endpoint via `HTTP`, making the `JAEGER_AGENT_HOST` and `JAEGER_AGENT_PORT` unused. If `JAEGER_ENDPOINT` is
secured, HTTP basic authentication can be performed by setting the `JAEGER_USER` and `JAEGER_PASSWORD` environment
variables.
