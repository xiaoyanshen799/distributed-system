package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/xiaoyanshen799/distributed-system/golang"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("pets.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移数据库模型
	if err := db.AutoMigrate(&PetModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterPetServiceServer(grpcServer, &server{db: db})

	// 注册反射服务，以便使用 grpcurl 等工具进行调试
	reflection.Register(grpcServer)

	log.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
