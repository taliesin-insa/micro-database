# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Emilie HUMMEL"

# Dépendances nécessaires pour compiler le fichier protocole
RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go get -u github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Define directory
ADD src /src
WORKDIR /src/mongoGo

# Download dependancies
RUN go get -u github.com/gorilla/mux
RUN go get -u go.mongodb.org/mongo-driver/bson
RUN go get -u go.mongodb.org/mongo-driver/mongo
RUN go get -u go.mongodb.org/mongo-driver/mongo/options

# Build all project
RUN go build .

# Command to run the executable
CMD ["./rest"]
