FROM golang

WORKDIR /app

COPY . .

RUN go build -o tcp-chat .

ENTRYPOINT [ "tcp-chat" ]