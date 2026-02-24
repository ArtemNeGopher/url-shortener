# ====================== Build Stage ======================
FROM golang:1.25.7-alpine AS builder

ARG SERVICE_PATH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags=prod -o /app/main ./services/analytics-service/cmd/

# ====================== Run Stage ======================
FROM alpine:3.23.3

ARG SERVICE_PATH

WORKDIR /app

COPY services/analytics-service/config/config.yaml ./config/
COPY services/analytics-service/migrations ./migrations
COPY --from=builder /app/main ./


EXPOSE 50053

CMD ["./main"]
