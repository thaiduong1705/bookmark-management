FROM golang:1.27.3-alpine AS test

RUN apk add --no-cache make bc gcc musl-dev
ENV CGO_ENABLED=1

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make test

FROM golang:1.27.3-alpine AS builder

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go tool swag init -g cmd/api/main.go
RUN go build -o bookmark_service cmd/api/main.go

FROM alpine:3.20.10 AS runner
WORKDIR /app

COPY --from=builder /opt/app/bookmark_service .

CMD ["./bookmark_service"]