# Этап, на котором выполняется сборка приложения
FROM golang:alpine as builder

COPY . .

RUN go mod download

RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go build -o /bin/main ./cmd/main.go

FROM alpine:latest
COPY --from=builder /bin/main /main
CMD ["./main"]