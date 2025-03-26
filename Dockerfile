FROM golang:1.23.7-alpine

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod vendor

CMD ["air", "-c", ".air.toml"]