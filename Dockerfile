# Use an official Go runtime as a parent image
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o test .

# Start a new stage for the minimal runtime container
FROM debian:buster-slim

# Set the working directory inside the minimal runtime container
WORKDIR /app

# Copy the built binary from the builder container into the minimal runtime container
COPY --from=builder /app/test .

# Expose the port your Go application listens on (if applicable)
# EXPOSE 8080

# Run your Go application
CMD ["./test"]