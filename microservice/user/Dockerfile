# syntax=docker/dockerfile:1

FROM golang:latest

LABEL maintainer="Artyom <artyomsh01@yandex.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY config .
COPY microservice .
COPY middleware .
COPY pkg .
COPY prometheus .
COPY server .
COPY service .
COPY main.go .

RUN cd microservice/user/cmd && go build -o user

EXPOSE 8084

WORKDIR /app/microservice/user/cmd

CMD ["./user"]
