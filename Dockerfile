FROM golang:latest

COPY ./ ./
RUN go build -o main .
CMD ["./main"]