FROM golang:1.12.17 AS builder
RUN mkdir /jbtklinebot
WORKDIR /jbtklinebot
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./

FROM alpine:3.10.3
RUN mkdir -p /jbtklinebot

COPY --from=builder /jbtklinebot/main /jbtklinebot/main
COPY --from=builder /jbtklinebot/app.yml /jbtklinebot/app.yml
COPY --from=builder /jbtklinebot/.env /jbtklinebot/.env
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /jbtklinebot
USER appuser

WORKDIR /jbtklinebot
CMD ./main

