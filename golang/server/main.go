package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "distributed-system"

	"github.com/google/uuid"
)

type server struct {
	pb.UnimplementedPetServiceServer
	mu sync.Mutex
	db *gorm.DB
}

type PetModel struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Gender  string
	Age     int32
	Breed   string
	Picture []byte
}

// BeforeCreate hook to set UUID before creating a new record
func (p *PetModel) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}

func (s *server) RegisterNewPet(_ context.Context, req *pb.RegisterNewPetRequest) (*pb.RegisterNewPetReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// get pet info from request and build db input
	pet := PetModel{
		Name:    req.GetName(),
		Gender:  req.GetGender(),
		Age:     req.GetAge(),
		Breed:   req.GetBreed(),
		Picture: req.GetPicture(),
	}

	// create db instance
	if err := s.db.Create(&pet).Error; err != nil {
		log.Printf("Error creating pet: %v", err)
		return &pb.RegisterNewPetReply{
			Code: 1,
			Msg:  "Failed to register pet",
		}, err
	}

	return &pb.RegisterNewPetReply{
		Code: 0,
		Msg:  "Pet registered successfully",
	}, nil
}

func (s *server) SearchPet(ctx context.Context, req *pb.SearchPetRequest) (*pb.SearchPetReply, error) {
	var pets []PetModel

	switch detail := req.GetDetail().(type) {
	case *pb.SearchPetRequest_Name:
		s.db.Where("name = ?", detail.Name).Find(&pets)
	case *pb.SearchPetRequest_Gender:
		s.db.Where("gender = ?", detail.Gender).Find(&pets)
	case *pb.SearchPetRequest_Age:
		s.db.Where("age = ?", detail.Age).Find(&pets)
	case *pb.SearchPetRequest_Breed:
		s.db.Where("breed = ?", detail.Breed).Find(&pets)
	default:
		return &pb.SearchPetReply{
			Pets: nil,
		}, nil
	}

	var responsePets []*pb.Pet
	for _, pet := range pets {
		responsePets = append(responsePets, &pb.Pet{
			Name:    pet.Name,
			Gender:  pet.Gender,
			Age:     pet.Age,
			Breed:   pet.Breed,
			Picture: pet.Picture,
		})
	}

	return &pb.SearchPetReply{
		Pets: responsePets,
	}, nil
}

func main() {
	// connect db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// the database table is created or updated to match the 'PetModel' struct.
	if err := db.AutoMigrate(&PetModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create gRPC server
	grpcServer := grpc.NewServer()

	// register server
	pb.RegisterPetServiceServer(grpcServer, &server{db: db})

	// register server
	reflection.Register(grpcServer)

	log.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
