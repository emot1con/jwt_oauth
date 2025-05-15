FROM golang:1.24-alpine

WORKDIR /app

# Copy frontend files
COPY . ./

# Install any dependencies if needed
RUN go mod init frontend-server && \
    go mod tidy

# Expose port 3000
EXPOSE 3000

# Run the server
CMD ["go", "run", "server.go"]
