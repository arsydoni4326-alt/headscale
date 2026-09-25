# Headscale server container image
#
# Builds the headscale control server from source and runs it in a minimal
# Debian runtime. The image is compatible with the official headscale
# container usage: `docker run ... headscale serve` (or `command: serve` in
# docker-compose).
#
# Build:
#   docker build --build-arg VERSION=$(git describe --always --tags --dirty) -t headscale .
#
# Run (mount your config and data directory):
#   docker run --rm -it \
#     -v "$(pwd)/config:/etc/headscale:ro" \
#     -v "$(pwd)/lib:/var/lib/headscale" \
#     -p 127.0.0.1:8080:8080 \
#     -p 127.0.0.1:9090:9090 \
#     headscale serve

# Build stage
FROM golang:1.27.0-trixie AS builder

ARG APP_VERSION=v0.0.0
ARG APP_COMMIT=unknown
ARG BUILD_DATE=2025-09-09
ENV GOPATH /go
WORKDIR /go/src/headscale

# Download dependencies first - only invalidated when go.mod/go.sum change
COPY go.mod go.sum /go/src/headscale/
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go build -buildmode=pie \
  -ldflags="-s -w -X 'main.Version=${CPA_VERSION}' -X 'main.Commit=${CPA_COMMIT}' -X 'main.BuildDate=${BUILD_DATE}'" \
  -o /go/bin/headscale ./cmd/headscale

# Runtime stage
FROM debian:trixie-slim

# ca-certificates is needed for outbound HTTPS (update checks, OIDC, ...)
RUN apt-get --update install --no-install-recommends --yes ca-certificates \
  && apt-get dist-clean

# Data and runtime directories used by the default configuration
RUN mkdir -p /var/lib/headscale /var/run/headscale

COPY --from=builder /go/bin/headscale /usr/local/bin/headscale

EXPOSE 8080/tcp 9090/tcp

ENTRYPOINT ["headscale"]
CMD ["serve"]