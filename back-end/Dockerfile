# Build stage
FROM golang:1.23-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o MegaMedia cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /root/
COPY --from=build /app/MegaMedia .
COPY --from=build /app/static ./static
EXPOSE 8080
CMD ["./MegaMedia"]
