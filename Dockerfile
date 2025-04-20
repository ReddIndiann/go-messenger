FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-messenger

EXPOSE 5000

CMD ["./go-messenger"] 