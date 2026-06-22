# syntax=docker/dockerfile:1

# Build stage
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache ca-certificates git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server
RUN make build-server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the server executable from the build stage
COPY --from=build /app/server .

# Expose the port the server runs on
EXPOSE 8080

# Run the server
CMD ["./server"]
