FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o frontend ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/frontend .
COPY templates ./templates
COPY style.css .
EXPOSE 8080
CMD ["./frontend"]
