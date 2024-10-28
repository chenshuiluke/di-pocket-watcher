# Start from the official Go image
FROM golang:1.23

# Install PostgreSQL client
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download
RUN go mod tidy

# Copy the source code
COPY . .

# Run sqlc generate
RUN sqlc generate

# Create a startup script
COPY ./scripts/start.sh /start.sh
RUN chmod +x /start.sh

# Command to run the application
CMD ["/start.sh"]