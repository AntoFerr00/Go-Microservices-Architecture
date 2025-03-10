# user-service/Dockerfile
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy top-level go.mod and go.sum from the root build context
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into /app
COPY . .

# Build only the user-service binary
RUN go build -o /bin/user-service ./user-service/main.go

# Final stage (small Alpine image)
FROM alpine:3.17
COPY --from=build /bin/user-service /bin/user-service
EXPOSE 8081
ENTRYPOINT ["/bin/user-service"]
