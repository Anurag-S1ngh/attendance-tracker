# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /server .

EXPOSE 8000

CMD ["/app/server"]
