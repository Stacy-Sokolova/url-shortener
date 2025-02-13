# Этап, на котором выполняется сборка приложения
FROM golang:alpine as builder
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o /bin/main ./cmd/main.go

FROM alpine:latest
COPY --from=builder /bin/main /main
CMD ["./main"]