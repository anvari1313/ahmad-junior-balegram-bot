# Build stage: compile a static, stripped binary.
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cache module downloads independently of source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/ahmad-junior-balegram-bot .

# Runtime stage: minimal Alpine with CA certificates for TLS to the Bale and Z.ai APIs.
FROM alpine:3.24

RUN apk add --no-cache ca-certificates && \
    adduser -D -H -u 10001 bot
USER bot

COPY --from=builder /out/ahmad-junior-balegram-bot /usr/local/bin/ahmad-junior-balegram-bot

ENTRYPOINT ["ahmad-junior-balegram-bot"]
