# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the container
COPY . .

# Download the Go module dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose the desired port
EXPOSE 8080

# Set the command to run the Go application
CMD ["./main"]
