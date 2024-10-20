# use Go official image
FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y \
    xdg-utils \
    curl \
    gcc \
    musl-dev \
    sqlite3 \
    libsqlite3-dev

WORKDIR /app

RUN mkdir -p /app/client_register
RUN mkdir -p /app/client_search
RUN mkdir -p /app/server
RUN mkdir /app/inputPicture

# copy all necessary file to /app directory
COPY ./client_register /app/client_register
COPY ./client_search /app/client_search
COPY ./server/main.go /app/server
COPY go.mod /app
COPY go.sum /app
COPY pet_grpc.pb.go /app/
COPY pet.pb.go /app/
COPY ./inputPicture /app/inputPicture

ENV GO111MODULE=on

RUN go mod tidy
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64
# build project
RUN go build -o server_out ./server/main.go
RUN go build -o register-client ./client_register/main.go
RUN go build -o search-client ./client_search/main.go
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev
