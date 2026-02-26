# ====================== Build Stage ======================
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GIN_MODE=release go build -tags=prod -o /app/main ./services/api-gateway/cmd/

# ====================== Run Stage ======================
FROM alpine:3.23.3

WORKDIR /app

COPY services/api-gateway/config/config.yaml ./config/
COPY services/api-gateway/frontend/index.html ./index.html
COPY --from=builder /app/main ./


EXPOSE 8080

CMD ["./main"]
