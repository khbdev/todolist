# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app

# wait-for-it skriptini qo'shamiz
COPY --from=builder /app/main .
COPY wait-for-it.sh .
COPY --from=builder /app/.env .

EXPOSE 8082

# CMD ni o'zgartiramiz: MySQL tayyor bo‘lgach ishga tushadi
CMD ["./wait-for-it.sh", "todolist-mysql:3306", "--timeout=30", "--strict", "--", "./main"]
