
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go


FROM alpine:3.18
WORKDIR /app


COPY --from=builder /app/main .

COPY --from=builder /app/.env .

EXPOSE 8082


CMD ["./main"]
