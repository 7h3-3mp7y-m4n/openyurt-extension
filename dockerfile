FROM golang:1.25rc3-alpine3.22 AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend
COPY ui ./ui
COPY scripts ./scripts
RUN chmod +x scripts/*.sh

RUN go build -o openyurt-backend ./backend/main.go
# why not save your space a little by multi-built image 🤗
FROM alpine:3.22.1

RUN apk add --no-cache bash curl ca-certificates

WORKDIR /app

COPY --from=builder /app/openyurt-backend .

COPY --from=builder /app/ui ./ui
COPY --from=builder /app/scripts ./scripts

RUN chmod +x scripts/*.sh

EXPOSE 8080

CMD ["./openyurt-backend"]