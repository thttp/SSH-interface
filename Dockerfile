FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o tui-portfolio .

# ── Imagem final ──────────────────────────────────────────────────────────────
FROM alpine:latest

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/tui-portfolio .

RUN mkdir -p .ssh

EXPOSE 2222

CMD ["./tui-portfolio", "--serve"]