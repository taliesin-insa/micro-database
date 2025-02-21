# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
# Bad practice but anyway
FROM golang:latest AS builder

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

# Define directory
ADD src /src
WORKDIR /src/micro-database

# Download dependancies (if you try to build your image without following lines you will see missing packages)
RUN go get -u github.com/gorilla/mux
RUN go get -u go.mongodb.org/mongo-driver/bson
RUN go get -u go.mongodb.org/mongo-driver/mongo
RUN go get -u go.mongodb.org/mongo-driver/mongo/options
RUN go get -u github.com/prometheus/client_golang/prometheus
RUN go get -u github.com/prometheus/client_golang/prometheus/promauto
RUN go get -u github.com/prometheus/client_golang/prometheus/promhttp
RUN go get -u github.com/taliesin-insa/lib-auth

# Build all project statically (prevent some exec user process caused "no such file or directory" error)
ENV CGO_ENABLED=0
RUN go build .

# Build the docker image from a lightest one (otherwise it weights more than 1Go)
FROM alpine:latest

# Expose port 8080 to the outside world
EXPOSE 8080

# Don't really know what this does
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy on the executive env
COPY --from=builder /src/micro-database/micro-database .

# Command to run the executable
CMD ["./micro-database"]
