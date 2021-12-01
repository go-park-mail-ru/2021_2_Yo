# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download


RUN go build 

EXPOSE 8080

CMD ["./backend"]