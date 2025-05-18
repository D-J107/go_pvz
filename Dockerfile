# Stage 1: Build
FROM golang:1.23-bullseye AS build

# Set the working directory inside the container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pvz-api ./cmd/pvz_application/main.go

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY --from=build /app/pvz-api .
EXPOSE 8080 3000 9000
USER 1001
CMD ["/pvz-api"]