import grpc
from concurrent import futures
import pet_pb2
import pet_pb2_grpc


class PetService(pet_pb2_grpc.PetServiceServicer):
    def __init__(self):
        self.pets = []

    def RegisterNewPet(self, request, context):
        new_pet = pet_pb2.Pet(
            name=request.name,
            gender=request.gender,
            age=request.age,
            breed=request.breed,
            picture=request.picture
        )
        self.pets.append(new_pet)
        return pet_pb2.RegisterNewPetReply(code=0, msg="Pet registered successfully")

    def SearchPet(self, request, context):
        results = []
        for pet in self.pets:
            if (request.HasField("name") and pet.name == request.name) or \
                    (request.HasField("gender") and pet.gender == request.gender) or \
                    (request.HasField("age") and pet.age == request.age) or \
                    (request.HasField("breed") and pet.breed == request.breed):
                results.append(pet)
        return pet_pb2.SearchPetReply(pets=results)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pet_pb2_grpc.add_PetServiceServicer_to_server(PetService(), server)
    server.add_insecure_port('[::]:50051')
    print("Server starting on port 50051...")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
