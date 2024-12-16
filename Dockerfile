FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o octo-manager ./cmd/server

# Create minimal runtime image - reduce size and improve security
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/octo-manager .

EXPOSE 8080

ENTRYPOINT ["./octo-manager"]