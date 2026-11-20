FROM golang:1.27.3-alpine AS test

RUN apk add --no-cache make bc

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make test

FROM golang:1.27.3-alpine AS builder

RUN mkdir -p "/opt/app"

WORKDIR /opt/app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
RUN swag init -g cmd/api/main.go
RUN go build -o bookmark_service cmd/api/main.go

FROM alpine:3.20.10 AS runner
WORKDIR /app

COPY --from=builder /opt/app/bookmark_service .

CMD ["./bookmark_service"]