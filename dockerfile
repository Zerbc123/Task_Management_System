FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

FROM alpine:latest AS server

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

FROM alpine:latest AS worker

WORKDIR /app

COPY --from=builder /app/worker .

CMD ["./worker"]