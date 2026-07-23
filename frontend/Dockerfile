# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY frontend/go.mod frontend/go.sum ./frontend/
RUN cd frontend && go mod download

COPY . .

WORKDIR /app/frontend
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /frontend ./cmd

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /frontend ./frontend

USER appuser

EXPOSE 8080

ENTRYPOINT ["./frontend"]