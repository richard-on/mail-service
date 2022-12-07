FROM golang:1.19.3-buster as builder

WORKDIR /mail

COPY go.* ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -ldflags "-X main.version=0.1.0 -X main.build=`date -u +.%Y%m%d.%H%M%S`" \
    -o run cmd/mail/main.go

FROM alpine:latest

WORKDIR /mail

COPY --from=builder /mail/run /mail/run
COPY --from=builder /mail/.env /mail/.env

EXPOSE 80

RUN mkdir -p /auth/logs && \
    apk update && apk add curl && apk add --no-cache bash && \
    apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ./run