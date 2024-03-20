FROM golang:1.22-alpine

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o unclatter cmd/api/main.go

EXPOSE 8080

CMD ["./unclatter"]
