FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o health-checker ./cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/health-checker /health-checker
EXPOSE 9090
ENTRYPOINT ["/health-checker"]
