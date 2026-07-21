# Stage 1: Build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o frontend ./cmd

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/frontend .
COPY templates ./templates
COPY static ./static

EXPOSE 8080
CMD ["./frontend"]