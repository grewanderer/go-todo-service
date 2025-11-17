FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-todo-service ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /go-todo-service /go-todo-service

EXPOSE 8080

ENV SERVER_PORT=8080

ENTRYPOINT ["/go-todo-service"]
