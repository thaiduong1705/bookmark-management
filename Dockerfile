FROM golang:1.27.3-alpine AS builder

RUN mkdir -p "/opt/app"

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark_service cmd/api/main.go

FROM alpine:3.20.10 AS runner

WORKDIR /app

COPY --from=builder /opt/app/bookmark_service .

CMD ["./bookmark_service"]