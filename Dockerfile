FROM golang:alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM alpine:latest

WORKDIR /app 

COPY --from=builder /server .
COPY config/ ./config/

RUN mkdir -p /app/data

EXPOSE 8080

ENTRYPOINT ["./server"]
CMD ["-port", "8080","-store","file","-bot-service-url","http://botservice:9090/move"]