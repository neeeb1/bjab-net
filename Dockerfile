ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /run-app

FROM scratch
WORKDIR /app
COPY --from=builder /run-app /app/run-app
COPY --from=builder /src/web /app/web
EXPOSE 3000
CMD ["/app/run-app"]