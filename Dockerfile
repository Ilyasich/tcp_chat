FROM golang AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tcp-chat .

# -------------------------------------------------
FROM alpine:edge

WORKDIR /app

ENV HELLO_MESSAGE="My new hello from docker"
COPY --from=builder /app/tcp-chat .

ENTRYPOINT [ "/app/tcp-chat" ]