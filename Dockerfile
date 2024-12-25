# Use an official Go image to ensure all dependencies are included
FROM golang:1.22.6

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the local project directory to the container
COPY . .

# Download Go dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose the application's port (8080)
EXPOSE 8080

# Run the application executable directly
CMD ["./main"]