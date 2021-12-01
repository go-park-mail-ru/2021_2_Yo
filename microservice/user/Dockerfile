# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN cd microservice/user/cmd && go build -o user

EXPOSE 8084

WORKDIR /app/microservice/user/cmd

CMD ["./user"]
