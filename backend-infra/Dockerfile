# Stage 1: Build
FROM golang:1.25.0-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/server ./cmd/server

# Stage 2: Run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /out/server /server
COPY --from=builder /app/docs ./docs
RUN apk add --no-cache ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

EXPOSE 3000
ENTRYPOINT ["/server"]

