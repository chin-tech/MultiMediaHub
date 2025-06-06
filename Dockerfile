FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY ./src .
RUN go mod download

RUN go build -o super_dict .

FROM alpine:latest

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/super_dict .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 9000

ENTRYPOINT ["./super_dict"]
