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

# 下载依赖
RUN go mod tidy

# 为 Linux ARM 架构构建服务器
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server ./server/main.go

# 为 Linux ARM 架构构建客户端
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o register-client ./client_register/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o search-client ./client_search/main.go

# 使用更小的基础镜像来运行应用程序
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 将构建的二进制文件复制到最终镜像
COPY --from=builder /app/server .
COPY --from=builder /app/register-client .
COPY --from=builder /app/search-client .

# 暴露 gRPC 服务器的端口
EXPOSE 50051

# 默认启动服务器
CMD ["./server"]

