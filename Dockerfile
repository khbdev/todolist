# =========================
# Builder stage
# =========================
FROM golang:1.24-alpine AS builder

# Kerakli paketlarni o‘rnatish (agar kerak bo‘lsa)
RUN apk add --no-cache git

WORKDIR /app

# Go mod fayllarini copy qilib, dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Faqat kodni copy qilish
COPY cmd ./cmd
COPY internal ./internal

# Build qilish
RUN go build -o main ./cmd/main.go

# =========================
# Final stage
# =========================
FROM alpine:3.18

WORKDIR /app

# Faqat build qilingan binary va .env faylini olib kelish
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Port
EXPOSE 8082

# CMD
CMD ["./main"]
