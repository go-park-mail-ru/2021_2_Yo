# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN cd cmd/auth && go build -o auth

EXPOSE 8081

WORKDIR /app/cmd/auth

CMD ["./auth"]
