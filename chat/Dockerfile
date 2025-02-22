# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY index.html .
EXPOSE 80
CMD ["./server"]
