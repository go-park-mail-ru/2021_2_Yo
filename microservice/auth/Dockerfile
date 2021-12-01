# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN cd microservice/auth/cmd && go build -o auth

EXPOSE 8081

WORKDIR /app/microservice/auth/cmd

CMD ["./auth"]
