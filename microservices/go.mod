module microservices

go 1.12

require (
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/echo-contrib v0.6.0
	github.com/labstack/echo/v4 v4.1.10
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber/jaeger-client-go v2.19.1-0.20191002155754-0be28c34dabf+incompatible
	golang.org/x/crypto v0.0.0-20191001170739-f9e2070545dc // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/labstack/echo-contrib => github.com/carlosedp/echo-contrib v0.6.1-0.20191002230647-bebaa759d659
