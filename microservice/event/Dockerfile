# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN cd microservice/event/cmd && go build -o event

EXPOSE 8083

WORKDIR /app/microservice/event/cmd

CMD ["./event"]
