# Start from the official Go image
FROM golang:1.22.5

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum Makefile ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY app/ ./app/
COPY example/ ./example/

# Build the application
RUN make build-example

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./bin/web-service-gin"]
