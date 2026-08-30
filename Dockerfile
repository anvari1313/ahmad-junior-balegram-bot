# Build stage: cross-compile a static, stripped binary for the target platform.
# Pinning the builder to $BUILDPLATFORM keeps Go compiling natively on the
# build machine; TARGETOS/TARGETARCH select the output, so multi-platform
# builds need no QEMU emulation.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Cache module downloads independently of source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" \
    -o /out/ahmad-junior-balegram-bot .

# Runtime stage: minimal Alpine. It contains only COPY instructions, so no
# target-arch process is ever executed and multi-platform builds stay
# emulation-free. The CA bundle (TLS to the Bale and Z.ai APIs) is taken from
# the builder image; the bot runs as an unprivileged user.
FROM alpine:3.24

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/ahmad-junior-balegram-bot /usr/local/bin/ahmad-junior-balegram-bot

USER 10001:10001

ENTRYPOINT ["ahmad-junior-balegram-bot"]
