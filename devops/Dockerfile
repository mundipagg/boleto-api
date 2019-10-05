# build stage
FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/mundipagg/boleto-api
COPY . .
RUN mkdir -p "deploy"

RUN go build -o /deploy/boleto-api

COPY devops/time "/deploy/lib/time/"
COPY boleto/favicon.ico "/deploy/boleto/"
COPY boleto/Arial.ttf "/deploy/boleto/"

FROM alpine:3.7
WORKDIR /home/mundipagg
COPY --from=build-env ./deploy/ .
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENTRYPOINT ["/home/mundipagg/boleto-api"]
EXPOSE 3000