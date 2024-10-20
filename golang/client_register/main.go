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

	pb "distributed-system"

	"google.golang.org/grpc"
)

func registerPet(client pb.PetServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	// input information
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

	// trans image to byte
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("unable to open picture: %v", err)
	}

	// create request
	pet := &pb.RegisterNewPetRequest{
		Name:    name,
		Gender:  gender,
		Age:     int32(age),
		Breed:   breed,
		Picture: imageData,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// call server function
	response, err := client.RegisterNewPet(ctx, pet)
	if err != nil {
		log.Fatalf("register pet error: %v", err)
	}
	log.Printf("register pet response: %s", response.GetMsg())
}

func main() {
	// Establish a gRPC connection to the 'server-container' at port 50051.
	// server-container is my server's container name
	conn, err := grpc.Dial("server-container:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client using the provided connection.
	client := pb.NewPetServiceClient(conn)

	registerPet(client)
}
