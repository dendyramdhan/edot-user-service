# Start from the official Go image to build your application
FROM golang:1.20-alpine as builder

# Create a directory for user service
WORKDIR /app

# Copy the go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build user service
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a Docker multi-stage build to create a lean production image
# Start from the alpine image to have a lightweight base
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
