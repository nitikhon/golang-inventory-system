FROM golang:1.23.7-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod vendor
RUN go build -o bin/server/main cmd/server/main.go
CMD ./bin/server/main