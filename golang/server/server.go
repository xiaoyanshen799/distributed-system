package main

import (
	"context"
	"log"

	pb "github.com/xiaoyanshen799/distributed-system"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedPetServiceServer
	db *gorm.DB
}

func (s *server) RegisterNewPet(_ context.Context, req *pb.RegisterNewPetRequest) (*pb.RegisterNewPetReply, error) {
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
