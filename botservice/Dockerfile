FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /botservice ./botservice

FROM alpine:latest

WORKDIR /app

COPY --from=builder /botservice .

EXPOSE 9090

ENTRYPOINT ["./botservice"]
CMD ["-botservice-port","9090"]