# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /expense-tracker ./cmd/server/main.go
# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /expense-tracker .
EXPOSE 8080
CMD ["./expense-tracker"]
