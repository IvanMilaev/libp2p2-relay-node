# Use Go Alpine image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the relay binary
RUN go build -o relay .

# Expose the port
EXPOSE 4001

# Run the relay node
CMD ["./relay"]

