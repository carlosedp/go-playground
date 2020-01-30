# This is the builder container
FROM carlosedp/golang:1.13 AS builder

WORKDIR $GOPATH/src/app
ADD . $GOPATH/src/app/
# RUN go get . # In case your application has dependencies
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o main .

RUN mv $GOPATH/src/app/main /

# Application container
FROM scratch

COPY --from=builder /main /main

EXPOSE 8080

ENTRYPOINT ["/main"]

