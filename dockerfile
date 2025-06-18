FROM golang:alpine AS build-env

RUN apk --no-cache add gcc musl-dev go-task-task

RUN mkdir /app
WORKDIR /app

RUN go install github.com/vektra/mockery/v3@v3.4.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN task gen

RUN go build -o surl ./cmds/main

FROM alpine:latest

RUN apk --no-cache add ca-certificates go-task-task postgresql-client

WORKDIR /app

COPY --from=build-env /app/Taskfile.yml ./
COPY --from=build-env /app/surl ./

EXPOSE 8080

CMD ["/app/surl", "http", "-p", "8080"]
