FROM golang:1.25.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/api/main.go

EXPOSE 8000


CMD ["./app"]