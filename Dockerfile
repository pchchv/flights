FROM golang:1.18-alpine

WORKDIR /app

COPY . /app

RUN go mod init main.go
RUN go mod tidy