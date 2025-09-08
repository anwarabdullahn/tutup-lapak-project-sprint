# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download || true
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

# Runtime stage
FROM alpine:3.20
WORKDIR /
COPY --from=builder /out/server /server
ENV PORT=3000
EXPOSE 3000
ENTRYPOINT ["/server"]

