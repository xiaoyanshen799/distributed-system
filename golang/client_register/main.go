package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/xiaoyanshen799/distributed-system/golang"

	"google.golang.org/grpc"
)

func registerPet(client pb.PetServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	// 获取宠物信息
	fmt.Print("please input pet's name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("please input pet's gender: ")
	gender, _ := reader.ReadString('\n')
	gender = strings.TrimSpace(gender)

	fmt.Print("please input pet's age: ")
	ageInput, _ := reader.ReadString('\n')
	age, err := strconv.Atoi(strings.TrimSpace(ageInput))
	if err != nil {
		log.Fatalf("invalid age: %v", err)
	}

	fmt.Print("please input pet's breed: ")
	breed, _ := reader.ReadString('\n')
	breed = strings.TrimSpace(breed)

	fmt.Print("please input picture path: ")
	imagePath, _ := reader.ReadString('\n')
	imagePath = strings.TrimSpace(imagePath)

	// 将图片读取为字节数组
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("unable to open picture: %v", err)
	}

	// 创建宠物请求
	pet := &pb.RegisterNewPetRequest{
		Name:    name,
		Gender:  gender,
		Age:     int32(age),
		Breed:   breed,
		Picture: imageData,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.RegisterNewPet(ctx, pet)
	if err != nil {
		log.Fatalf("register pet error: %v", err)
	}
	log.Printf("register pet response: %s", response.GetMsg())
}

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPetServiceClient(conn)

	// 调用注册宠物的方法
	registerPet(client)
}
