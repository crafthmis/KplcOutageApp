# Start from the official Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .
COPY config/env.json /app/config/env.json
COPY backup.sql /docker-entrypoint-initdb.d/backup.sql

# Build the application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]