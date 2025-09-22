# Gunakan base image Go alpine
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy dependency files dulu untuk caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary supaya lebih cepat di production (optional)
RUN go build -o app main.go

# Expose port container (sesuai CMD di docker-compose)
EXPOSE 8080

# Jalankan binary
CMD ["./app"]
