# order-service/Dockerfile
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/order-service ./order-service/main.go

FROM alpine:3.17
COPY --from=build /bin/order-service /bin/order-service
EXPOSE 8082
ENTRYPOINT ["/bin/order-service"]
