FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o secure-core ./cmd/server

# ---------- runtime ----------
FROM alpine:3.20

WORKDIR /app

RUN adduser -D secure
USER secure

COPY --from=builder /app/secure-core /app/secure-core
COPY configs /app/configs

EXPOSE 4000
EXPOSE 8080

CMD ["./secure-core"]
