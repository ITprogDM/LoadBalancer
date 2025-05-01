FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git

# Копируем go.mod и go.sum отдельно — это кэшируется
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Путь до main.go должен быть точным
RUN go build -o loadbalancer ./cmd/main.go

# Проверим, что бинарник точно существует
RUN ls -la /app/loadbalancer

EXPOSE 8080

CMD ["/app/loadbalancer"]