# Stage 1: Build
FROM golang:1.21 AS builder

# Install system dependencies
RUN apt-get update && apt-get install -y git libsensors4-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project from the subdirectory
COPY sensor-exporter/ ./sensor-exporter/

# Move into the Go project directory
WORKDIR /app/sensor-exporter

# Download dependencies and generate go.sum
RUN go mod tidy

# Build the binary
RUN go build -o /sensor-exporter

RUN ls

# Stage 2: Minimal runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache lm_sensors

# Copy built binary from the builder stage
COPY --from=builder /sensor-exporter /usr/local/bin/sensor-exporter

EXPOSE 9255

ENTRYPOINT ["/usr/local/bin/sensor-exporter"]
CMD ["--web.listen-address=:9255"]

