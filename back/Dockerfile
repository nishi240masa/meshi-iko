# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS base
WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS dev
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/air-verse/air@v1.62.0 \
    && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]


FROM base AS builder
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api


FROM alpine:3.22 AS prod
RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/api /app/api

RUN mkdir -p /app/data && chown -R app:app /app
USER app
ENV GIN_MODE=release \
    PORT=8080 \
    DB_PATH=/app/data/meshi-iko.db \
    TZ=Asia/Tokyo
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/v1/health || exit 1
ENTRYPOINT ["/app/api"]
