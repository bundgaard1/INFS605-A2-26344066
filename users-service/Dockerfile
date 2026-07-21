# ==========================================
# STAGE 1: Build phase
# ==========================================
FROM golang:1.26-alpine AS builder

# Installer ca-certificates og git (hvis der skal hentes eksterne pakker over HTTPS)
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# 1. Kopier afhængighedsfiler først (optimering af Docker cache layer)
COPY go.mod go.sum ./
RUN go mod download

# 2. Kopier kildekoden (inklusiv den genererede /gen mappe)
COPY . .

# 3. Kompiler Go-binærfilen statisk uden CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /users-service ./cmd/main.go

# ==========================================
# STAGE 2: Minimal Runtime phase
# ==========================================
FROM alpine:3.20

# Tilføj CA certifikater og opret en non-root bruger af sikkerhedshensyn
RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Kopier den kompilerede binære fil fra builder-staget
COPY --from=builder /users-service .

# Skift til non-root bruger
USER appuser

# Eksponer gRPC-porten
EXPOSE 50051

# Kør servicen
ENTRYPOINT ["./users-service"]