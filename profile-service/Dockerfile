# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git 

WORKDIR /app

COPY profile-service/go.mod profile-service/go.sum ./profile-service/
RUN cd profile-service && go mod download

COPY . .

# 3. Skift til profile-service mappen og kompiler binærfilen
WORKDIR /app/profile-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /profile-service ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

RUN chown -R appuser:appgroup /app

COPY --from=builder /profile-service .

USER appuser

EXPOSE 50051

ENTRYPOINT ["./profile-service"]