FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main ./src/main.go

FROM golang:1.17

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
