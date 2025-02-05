# Build stage
FROM golang:1.23-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o MegaMedia cmd/main.go

# Run stage
FROM alpine:3.18
ENV GO_ENV=production
WORKDIR /root/
COPY --from=build /app/MegaMedia .
COPY --from=build /app/static ./static
COPY .env.production .env.development ./
EXPOSE 8080
CMD ["./MegaMedia"]
