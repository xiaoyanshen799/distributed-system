import grpc
import pet_pb2
import pet_pb2_grpc


def register_pet(stub):
    # Simulating pet data for registration
    # pet = pet_pb2.Pet(
    #     name="Bella",
    #     gender="Female",
    #     age=3,
    #     breed="Labrador",
    #     picture=b"picture_bytes_here"  # Simulating a picture as bytes
    # )
    # response = stub.RegisterNewPet(pet)
    request = pet_pb2.RegisterNewPetRequest(
        name="Bella",
        gender="Female",
        age=3,
        breed="Labrador",
        picture="picture_bytes_here"  # This should be a string, not bytes
    )
    response = stub.RegisterNewPet(request)
    print(f"Register Response: {response.msg}")


def search_pet_by_name(stub, name):
    # Searching for pet by name
    search_request = pet_pb2.SearchPet(name=name)
    response = stub.SearchPet(search_request)
    if response.pets:
        for pet in response.pets:
            print(f"Found Pet: {pet.name}, {pet.gender}, {pet.age} years, Breed: {pet.breed}")
    else:
        print("No pets found.")


def search_pet_by_breed(stub, breed):
    # Searching for pet by breed
    search_request = pet_pb2.SearchPetRequest(breed=breed)
    response = stub.SearchPet(search_request)
    if response.pets:
        for pet in response.pets:
            print(f"Found Pet: {pet.name}, {pet.gender}, {pet.age} years, Breed: {pet.breed}")
    else:
        print("No pets found.")


def run():
    # Connect to the gRPC server
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = pet_pb2_grpc.PetServiceStub(channel)

        # Register a pet
        register_pet(stub)

        # Search for pets
        search_response = stub.SearchPet(pet_pb2.SearchPetRequest(breed="Labrador"))
        print("Search results:")
        for pet in search_response.pets:
            print(f"Name: {pet.name}, Breed: {pet.breed}, Age: {pet.age}")

        # # Search for pets by name
        # search_pet_by_name(stub, "Bella")
        #
        # # Search for pets by breed
        # search_pet_by_breed(stub, "Labrador")


if __name__ == "__main__":
    run()