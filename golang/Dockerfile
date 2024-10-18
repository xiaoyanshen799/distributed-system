# 使用 Go 的官方镜像作为构建环境
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

RUN mkdir -p /app/client_register
RUN mkdir -p /app/client_search
RUN mkdir -p /app/server

# 将当前目录的所有文件复制到容器的 /app 目录
COPY ./client_register /app/client_register
COPY ./client_search /app/client_search
COPY ./server/main.go /app/server
COPY go.mod /app
COPY go.sum /app

ENV GO111MODULE=on
