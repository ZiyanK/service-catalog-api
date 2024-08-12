FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service-catalog app/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8080

CMD ["./service-catalog"]