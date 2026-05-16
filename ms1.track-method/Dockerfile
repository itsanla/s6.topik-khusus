# ============================================================
# Stage 1: Build
# ============================================================
FROM golang:1.22-alpine AS builder

# Install dependensi build yang diperlukan
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk memanfaatkan layer cache Docker
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary dengan optimasi ukuran
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/track-method \
    ./main.go

# ============================================================
# Stage 2: Runtime (Alpine ringan)
# ============================================================
FROM alpine:3.20

# Install sertifikat SSL dan timezone data
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/track-method .

# Jalankan sebagai non-root user demi keamanan
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./track-method"]
