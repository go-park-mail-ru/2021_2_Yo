# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN cd cmd/user && go build -o user

EXPOSE 8084

WORKDIR /app/cmd/user

CMD ["./user"]
