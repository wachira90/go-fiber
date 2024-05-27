# Define the base image for building the Go application
FROM public.ecr.aws/docker/library/golang:1.22.1-alpine3.19 AS builder

# Set the working directory for the build stage
WORKDIR /app

# Copy your Go source code
COPY main.go .

# Install dependencies
RUN go mod init app
RUN go mod tidy

# Build the application binary
RUN CGO_ENABLED=0 go build -o main .

# Define the final image for running the application
FROM public.ecr.aws/docker/library/alpine:3.18.6

# Set the working directory for the final image
WORKDIR /app

# Copy only the application binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application listens on (optional)
EXPOSE 8080

# Specify the command to run the application
CMD ["./main"]
