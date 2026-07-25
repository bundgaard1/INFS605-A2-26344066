# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git 

WORKDIR /app

COPY notification-service/go.mod notification-service/go.sum ./notification-service/
RUN cd notification-service && go mod download

COPY . .

# 3. Skift til notification-service mappen og kompiler binærfilen
WORKDIR /app/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /notification-service ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

RUN chown -R appuser:appgroup /app

COPY --from=builder /notification-service .

USER appuser

EXPOSE 50052

ENTRYPOINT ["./notification-service"]