FROM golang:alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /frontend ./web

FROM alpine:latest

WORKDIR /app 

COPY --from=builder /frontend .

COPY web/*.html web/*.css web/*.js ./

EXPOSE 3000

ENTRYPOINT ["./frontend"]
CMD ["-port", "3000"]