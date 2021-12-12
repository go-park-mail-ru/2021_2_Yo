# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN cd cmd/server && go build -o server

EXPOSE 8080

WORKDIR /app/cmd/server

CMD ["./server"]