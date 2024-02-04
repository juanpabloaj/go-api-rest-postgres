FROM golang:1.21-alpine3.19 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
