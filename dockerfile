FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/my-app

RUN go mod tidy
RUN go build -o /app/api-go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-go /app/api-go

CMD ["./api-go"]