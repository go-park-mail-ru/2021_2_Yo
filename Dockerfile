# syntax=docker/dockerfile:1

FROM golang:latest as build
LABEL maintainer="Artyom <artyomsh01@yandex.ru>"
WORKDIR /app
RUN apt-get update && apt-get install -y \
    libwebp-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN make build

FROM golang:latest as server-build
LABEL maintainer="Artyom <artyomsh01@yandex.ru>"
WORKDIR /app
RUN apt-get update && apt-get install -y \
    libwebp-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN make build server
WORKDIR /app/bin/api
CMD ["./server"]

FROM golang:latest as auth-build
LABEL maintainer="Artyom <artyomsh01@yandex.ru>"
WORKDIR /app
RUN apt-get update && apt-get install -y \
    libwebp-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN make build auth-service
WORKDIR /app/bin/auth-service
CMD ["./auth"]

FROM golang:latest as event-build
LABEL maintainer="Artyom <artyomsh01@yandex.ru>"
WORKDIR /app
RUN apt-get update && apt-get install -y \
    libwebp-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN make build event-service
WORKDIR /app/bin/event-service
CMD ["./event"]

FROM golang:latest as user-build
LABEL maintainer="Artyom <artyomsh01@yandex.ru>"
WORKDIR /app
RUN apt-get update && apt-get install -y \
    libwebp-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN make build user-service
WORKDIR /app/bin/user-service
CMD ["./user"]