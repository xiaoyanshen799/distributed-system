package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	pb "github.com/xiaoyanshen799/distributed-system"

	"google.golang.org/grpc"
)

func saveImage(pet string, imageData []byte) (string, error) {
	// 创建图片保存目录
	imageDir := "downloaded_images"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		os.Mkdir(imageDir, os.ModePerm)
	}

	// 保存图片
	imagePath := filepath.Join(imageDir, pet+".jpg")
	err := os.WriteFile(imagePath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("unable to save picture: %v", err)
	}

	fmt.Printf("save image: %s\n", imagePath)
	return imagePath, nil
}

func openImage(imagePath string) error {
	// 根据操作系统运行不同的命令
	var cmd *exec.Cmd
	switch os := runtime.GOOS; os {
	case "darwin": // macOS
		cmd = exec.Command("open", imagePath)
	case "linux": // Linux
		cmd = exec.Command("xdg-open", imagePath)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", imagePath)
	default:
		return fmt.Errorf("不支持的操作系统: %s", os)
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("无法打开图片: %v", err)
	}

	return nil
}

func searchPet(client pb.PetServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("please input search keyword like: name:Labrador or gender:Female:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 解析用户输入的关键字
	splitInput := strings.Split(input, ":")
	if len(splitInput) != 2 {
		log.Fatalf("invalid type, please use this format, for example: name:Labrador")
	}
	field := splitInput[0]
	keyword := splitInput[1]

	var searchRequest *pb.SearchPetRequest

	// 根据输入的关键字构造对应的搜索请求
	switch field {
	case "name":
		searchRequest = &pb.SearchPetRequest{
			Detail: &pb.SearchPetRequest_Name{ // 使用包装器类型
				Name: keyword,
			},
		}
	case "gender":
		searchRequest = &pb.SearchPetRequest{
			Detail: &pb.SearchPetRequest_Gender{ // 使用包装器类型
				Gender: keyword,
			},
		}
	case "age":
		age, err := strconv.Atoi(keyword)
		if err != nil {
			log.Fatalf("无效的年龄: %v", err)
		}
		searchRequest = &pb.SearchPetRequest{
			Detail: &pb.SearchPetRequest_Age{ // 使用包装器类型
				Age: int32(age),
			},
		}
	case "breed":
		searchRequest = &pb.SearchPetRequest{
			Detail: &pb.SearchPetRequest_Breed{ // 使用包装器类型
				Breed: keyword,
			},
		}
	default:
		log.Fatalf("invalid field: %s", field)
	}

	// 调用 SearchPet 方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SearchPet(ctx, searchRequest)
	if err != nil {
		log.Fatalf("unable to search: %v", err)
	}

	// 打印搜索结果
	log.Printf("find %d pets:", len(resp.Pets))
	for _, pet := range resp.Pets {
		log.Printf("name: %s, gender: %s, age: %d, breed: %s", pet.Name, pet.Gender, pet.Age, pet.Breed)
		imagePath, err := saveImage(fmt.Sprintf("%s, %s, %d, %s", pet.Name, pet.Gender, pet.Age, pet.Breed), pet.Picture)
		if err != nil {
			log.Printf("无法保存图片: %v", err)
		}

		// 打开图片
		err = openImage(imagePath)
		if err != nil {
			log.Printf("无法打开图片: %v", err)
		}
		time.Sleep(2000)
	}
}

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPetServiceClient(conn)

	// 调用搜索宠物的方法
	searchPet(client)
}
