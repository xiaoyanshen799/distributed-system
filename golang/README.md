
# Distributed System Project

This project demonstrates a distributed system with components for client registration, client search, and a server. The project includes communication between these components using gRPC, SQLite for database storage, and Go for the backend logic.

## How to Compile and Run the Program

### Prerequisites

1. **Go**: Make sure you have Go installed on your machine (at least version 1.23).
2. **Docker**: Docker should be installed if you want to run the project using Docker containers.
3. **gRPC**: The project uses gRPC for communication between components, so ensure that the necessary gRPC tools and dependencies are installed.

### Running the Project Using Docker

The project includes a Dockerfile to help build and run the server and client components inside containers.

#### 0. Clone the repository: ####

   ```bash
   git clone https://github.com/xiaoyanshen799/distributed-system.git
   cd distributed-system/golang
   ```


#### 1. Create a Docker network

First, create a Docker network that the containers will use to communicate with each other:

```bash
docker network create my-network
```

This command creates a Docker network named `my-network`, allowing the server and client containers to communicate.

#### 2. Build the Docker image

Build the Docker image from the Dockerfile provided in the project:

```bash
docker build -t your_image_name .
```

Replace `your_image_name` with the name you'd like to give your Docker image.

#### 3. Run the server in a container

Now that the image is built, run the server in its own container:

```bash
docker run -d --name server-container --network my-network -p 50051:50051 your_image_name ./server_out
```

This command does the following:
- Runs the server in detached mode (`-d`).
- Assigns the container the name `server-container`.
- Connects the container to the `my-network` Docker network.
- Maps port 50051 on the host to port 50051 inside the container.
- Executes the `./server_out` binary (the server application) inside the container.

#### 4. Run the client in another container

Next, run a client container in interactive mode:

```bash
docker run -it --rm --network my-network your_image_name
```

This will start the client container and connect it to the same `my-network`, allowing it to communicate with the server.

Once inside the container, manually run the client applications:

```bash
# Inside the client container, run the client registration
./register-client

# Alternatively, run the client search
./search-client
```

You can run either the `register-client` or `search-client` as needed, and both will communicate with the server running in the `server-container`.

#### the picture use for register is in inputPicture/  ####
#### since unable to open picture in container, the picture of result in search-client will store in /downloaded_images ####
### Compiling the Project Manually

If you want to compile and run the project directly on your machine without Docker:

1. Clone the repository:

   ```bash
   git clone https://github.com/xiaoyanshen799/distributed-system.git
   cd distributed-system/golang
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Modify client's code:

   ```go
   //you should change the server-container into your real ip address in client_search/main.go and client_register/main.go in grpc.Dial
   conn, err := grpc.Dial("server-container:50051", grpc.WithInsecure(), grpc.WithBlock())
   ```

4. Compile the project:

   ```bash
   # Compile the server
   go build -o server-out ./server/main.go

   # Compile the client components
   go build -o register-client ./client_register/main.go
   go build -o search-client ./client_search/main.go
   ```

5. Run the server:

   ```bash
   ./server-out
   ```

6. Run the clients in another terminal:

   ```bash
   ./register-client
   ./search-client
   ```

## Anything Unusual About the Solution

- **CGO and SQLite**: The solution uses the `go-sqlite3` library, which requires CGO to be enabled. The Dockerfile is designed to enable CGO during the build process and install necessary dependencies like `gcc` and `sqlite3` development libraries.
  
- **xdg-open in Docker**: Since this project involves using `xdg-open`, which is typically available in graphical environments, but Docker containers do not usually include graphical environments, some functionalities like opening files in the default application may not work inside the container. If necessary, you can copy files to the host machine and view them locally.

- **Image Handling**: The project includes image files in `downloaded_images` that are used by the `client_search` component. These images need to be processed or handled in a headless environment due to the limitations of Docker’s containerized environment.

## External Sources Referenced

- The official [Go gRPC documentation](https://grpc.io/docs/languages/go/) was referenced for the setup of gRPC communication between the client and server components.
- [Go SQLite documentation](https://github.com/mattn/go-sqlite3) was used to understand how to integrate SQLite with Go and CGO.
- Official Docker guide for Go projects was referenced from [Docker Go Guide](https://docs.docker.com/guides/golang/).
- Protocol Buffers was set up according to the Proto3 documentation found at [Protobuf Documentation](https://protobuf.dev/programming-guides/proto3/).
