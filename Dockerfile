FROM golang:1.25-alpine 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate
RUN  go build -o ./server ./cmd/server/main.go

EXPOSE 8000

CMD ["/app/server"]
