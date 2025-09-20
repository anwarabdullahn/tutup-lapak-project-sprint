FROM golang:1.24 as builder
WORKDIR /src
COPY go.mod .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /server /server
EXPOSE 3003
USER 65532:65532
ENTRYPOINT ["/server"]
