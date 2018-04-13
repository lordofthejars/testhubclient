FROM golang:1.9.2-alpine3.7 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git glide

WORKDIR /go/src/github.com/lordofthejars/testhubclient

COPY . .

RUN glide install
RUN GOOS=linux GOARCH=amd64 go build -o binaries/testhubclient

FROM alpine:3.7
RUN addgroup -S testhub && adduser -S -G testhub testhub
USER testhub

COPY --from=builder /go/src/github.com/lordofthejars/testhubclient/binaries/testhubclient /usr/local/bin/testhubclient

ENTRYPOINT ["/usr/local/bin/testhubclient"]
CMD ["--help"]
