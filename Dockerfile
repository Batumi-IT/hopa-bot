FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hopa

FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /app/hopa /hopa
WORKDIR /app
CMD ["/hopa"]
