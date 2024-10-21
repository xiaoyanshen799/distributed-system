
# Pet Management System

This repository contains a simple gRPC-based pet management system, with both client and server components. The system allows for registering new pets and searching for pets based on various attributes (name, gender, age, breed).

## Project Structure

```
├── Dockerfile                # Docker image configuration for the project
├── client.py                 # Client-side gRPC code to interact with the server
├── docker-compose.yml         # Docker Compose file for running the project with dependencies
├── gui.py                    # A GUI interface for the pet management system
├── pet_pb2.py                # Generated code for protocol buffer message classes
├── pet_pb2_grpc.py           # Generated code for gRPC service classes
├── requirements.txt          # Python dependencies for the project
└── server.py                 # Server-side gRPC code to handle pet registration and searching
```

## Server (server.py)

The server provides two main functionalities:
1. **RegisterNewPet**: Registers a new pet with attributes like name, gender, age, breed, and picture.
2. **SearchPet**: Searches for pets based on the specified attributes (name, gender, age, breed).

### How to Run the Server

1. Install dependencies:
    ```bash
    pip install -r requirements.txt
    ```

2. Start the gRPC server:
    ```bash
    python server.py
    ```

   The server will start on port `50051`.

## Client (client.py)

The client is used to interact with the gRPC server to register pets and perform search operations.

### How to Run the Client

1. Install dependencies:
    ```bash
    pip install -r requirements.txt
    ```

2. Run the client:
    ```bash
    python client.py
    ```

## GUI (gui.py)

A graphical user interface (GUI) that interacts with the pet management system, allowing users to easily register and search for pets.

### Required Libraries
Ensure that `python-tk` is installed for the GUI to work correctly. On Ubuntu, you can install it with the following command:
```bash
sudo apt-get install python3-tk
```

For other operating systems, check the equivalent package manager commands.

## Docker

To run the entire system using Docker, build the image and run it with Docker Compose:

1. Build the Docker image:
    ```bash
    docker build -t pet-management .
    ```

2. Run the application with Docker Compose:
    ```bash
    docker-compose up
    ```

## Protocol Buffers

The `pet_pb2.py` and `pet_pb2_grpc.py` files contain the generated code from the `.proto` file, which defines the gRPC services and message types. You can regenerate these files by using the following command if changes are made to the `.proto` file:

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. path/to/your/pet.proto
```

## Conda Environment

It is recommended to run this project inside a Conda environment. You can create and activate an environment with the following commands:

```bash
conda create -n pet-management python=3.9
conda activate pet-management
```

## Requirements

Make sure the following dependencies are installed before running the project:

- Python 3.x
- gRPC (`grpcio`, `grpcio-tools`)
- `python-tk` for GUI support
- Other dependencies listed in `requirements.txt`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
