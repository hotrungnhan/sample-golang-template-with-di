FROM golang:alpine AS build-env

RUN apk --no-cache add gcc musl-dev go-task-task

RUN mkdir /app
WORKDIR /app

COPY . /app

RUN go mod download

RUN go install github.com/vektra/mockery/v3@v3.4.0

RUN task gen

RUN go build -o surl ./cmds/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build-env /app/app .


EXPOSE 8080

CMD ["surl http --port 8080"]
