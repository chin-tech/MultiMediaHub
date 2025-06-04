FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o mediahub .

FROM alpine:latest

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/mediahub .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 9000

ENTRYPOINT ["./mediahub"]
