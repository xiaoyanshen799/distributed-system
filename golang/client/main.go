package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/xiaoyanshen799/distributed-system"
	"google.golang.org/grpc"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPetServiceClient(conn)

	// 注册新宠物
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registerResp, err := client.RegisterNewPet(ctx, &pb.RegisterNewPetRequest{
		Name:    "Buddy",
		Gender:  "Male",
		Age:     3,
		Breed:   "Golden Retriever",
		Picture: "pictureData",
	})
	if err != nil {
		log.Fatalf("could not register pet: %v", err)
	}

	fmt.Printf("Pet registered code: %d, Message: %s\n", registerResp.GetCode(), registerResp.GetMsg())

	// 搜索宠物
	searchResp, err := client.SearchPet(ctx, &pb.SearchPetRequest{
		Detail: &pb.SearchPetRequest_Name{
			Name: "Buddy",
		},
	})
	if err != nil {
		log.Fatalf("could not search pet: %v", err)
	}

	fmt.Println("Search Results:")
	for _, pet := range searchResp.GetPets() {
		fmt.Printf("Name: %s, Gender: %s, Age: %d, Breed: %s\n",
			pet.GetName(),
			pet.GetGender(),
			pet.GetAge(),
			pet.GetBreed(),
		)
	}
}
