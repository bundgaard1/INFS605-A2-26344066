# ==========================================
# STAGE 1: Build phase (Monorepo Context)
# ==========================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

# 1. Optimering af Docker Cache: Hent afhængigheder først
COPY frontend/go.mod frontend/go.sum ./frontend/
RUN cd frontend && go mod download

# 2. Kopier kildekoden fra monorepoet
COPY . .

# 3. Kompiler Go-binærfilen statisk til Alpine Linux
WORKDIR /app/frontend
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /frontend ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

# Tilføj CA-certifikater (nødvendigt for udstrakte HTTPS-kald) og opret non-root bruger
RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /frontend .
# Kopier templates-mappen direkte ind under /app/templates
COPY --from=builder /app/frontend/templates ./templates
COPY --from=builder /app/frontend/static ./static

# Kopier den statisk kompilerede binærfil
COPY --from=builder /frontend .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./frontend"]