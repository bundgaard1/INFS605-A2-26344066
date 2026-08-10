# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git 

WORKDIR /app

COPY course-content-service/go.mod course-content-service/go.sum ./course-content-service/
RUN cd course-content-service && go mod download

COPY . .

WORKDIR /app/course-content-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /course-content-service ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

RUN chown -R appuser:appgroup /app

COPY --from=builder /course-content-service .

USER appuser

EXPOSE 50053

ENTRYPOINT ["./course-content-service"]