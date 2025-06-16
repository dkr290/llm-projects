# Use the official Golang image as a base
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o desktop-automation ./cmd/desktop-automation

# Command to run the executable
CMD ["./desktop-automation"]