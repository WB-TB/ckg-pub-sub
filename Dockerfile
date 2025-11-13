## DEVELOPMENT STAGE
FROM golang:1.25-trixie AS development

# Set direktori kerja
WORKDIR /app

# Copy file go.mod dan go.sum untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy sisa source code
COPY . .

# Install Watch tool untuk live reload saat development
RUN go install github.com/air-verse/air@latest

# Instal Delve untuk debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Build Producer
RUN go build -o /app/producer /app/cmd/producer/main.go

EXPOSE 2345

# CMD ["air", "--build.cmd", "go build -o /app/consumer /app/cmd/consumer/main.go", "--build.bin", "/app/consumer", "--debug.host", "0.0.0.0", "--debug.port", "2345"]
CMD ["air"]

## BUILDER STAGE
FROM golang:1.25-trixie AS builder

# Set direktori kerja
WORKDIR /app

# Copy file go.mod dan go.sum untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy sisa source code
COPY . .

# Build aplikasi Go. CGO_ENABLED=0 untuk static binary
RUN CGO_ENABLED=0 GOOS=linux -o /app/consumer /app/cmd/consumer/main.go
RUN CGO_ENABLED=0 GOOS=linux -o /app/producer /app/cmd/producer/main.go

## PRODUCTION STAGE
FROM debian:trixie-slim AS production

WORKDIR /root/
COPY --from=builder /app/consumer .
COPY --from=builder /app/producer .
CMD ["./consumer"]