# ==============================
# 1) BUILD STAGE
# ==============================
FROM golang:1.23-alpine AS builder

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o main \
    ./cmd/citizen-cab-svc

# ==============================
# 2) FINAL STAGE
# ==============================
FROM alpine:latest

WORKDIR /service

RUN apk add --no-cache ca-certificates

# 2.2. Копируем только Go-бинарь и конфиг
COPY --from=builder /service/main .
COPY config_docker.yaml .

ENV KV_VIPER_FILE=/service/config_docker.yaml
EXPOSE 8002

CMD ["./main", "run", "service"]
