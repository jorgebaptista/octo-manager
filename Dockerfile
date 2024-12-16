FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Build static binary
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o octo-manager ./cmd/server

# Create minimal runtime image - reduce size and improve security
FROM gcr.io/distroless/static

WORKDIR /app
COPY --from=builder /app/octo-manager .

EXPOSE 8080
ENTRYPOINT ["./octo-manager"]