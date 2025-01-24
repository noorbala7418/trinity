FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

FROM alpine:3.20

RUN apk add --no-cache tzdata

COPY --from=builder /app/main /app/trinity

CMD ["/app/trinity"]
