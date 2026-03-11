FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/pizza-hub ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && addgroup -g 1001 appgroup && adduser -D -u 1001 -G appgroup appuser

WORKDIR /app

COPY --from=builder /build/pizza-hub .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./pizza-hub"]

