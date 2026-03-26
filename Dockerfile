FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o quickslot ./cmd/app

# ---
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/quickslot .
COPY --from=builder /app/database/migrations ./database/migrations

EXPOSE 8080

CMD ["./quickslot"]
