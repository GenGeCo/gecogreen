# ============================================
# GECOGREEN - Backend Dockerfile
# ============================================

FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# ==========================================
# Stage 2: Runtime
# ==========================================
FROM alpine:3.19

RUN addgroup -g 1001 -S gecogreen && \
    adduser -u 1001 -S gecogreen -G gecogreen

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R gecogreen:gecogreen /app

USER gecogreen

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1

ENTRYPOINT ["/app/server"]
