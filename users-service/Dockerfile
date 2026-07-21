# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

# 1. Kopier kun modulfilerne for users-service for optimal Docker-caching
COPY users-service/go.mod users-service/go.sum ./users-service/
RUN cd users-service && go mod download

# 2. Kopier hele monorepoet (giver adgang til users-service og eventuelle fælles filer)
COPY . .

# 3. Skift til users-service mappen og kompiler binærfilen
WORKDIR /app/users-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /users-service ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Kopier den kompilerede binære fil fra builder-staget
COPY --from=builder /users-service .

USER appuser

EXPOSE 50051

ENTRYPOINT ["./users-service"]