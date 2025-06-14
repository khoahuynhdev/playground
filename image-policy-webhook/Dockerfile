# Use a slim version of Go >1.22.5 as the base image
FROM golang:1.24 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the application with optimization flags
RUN go build -ldflags="-s -w" -o image-policy-webhook main.go

# Use a minimal base image for the final container
FROM golang:1.24 AS final

# Create and use a non-root user
RUN useradd backend
USER backend

# Set the working directory
WORKDIR /home/backend

# Copy the built binary from the builder stage
COPY --from=builder /app/image-policy-webhook .

# Expose the port the application will run on
EXPOSE 8081

# Command to run the application
CMD ["./image-policy-webhook"]
