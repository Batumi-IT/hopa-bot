FROM golang:1.26.2-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hopa

FROM gcr.io/distroless/static-debian13:latest
COPY --from=builder /app/hopa /hopa
WORKDIR /app
CMD ["/hopa"]
