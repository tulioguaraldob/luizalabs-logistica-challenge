FROM golang:1.24.1

WORKDIR /app

COPY go.mod go.sum ./

COPY cmd /app

COPY . .

RUN go build -o cmd/main .

EXPOSE 8080

CMD ["./cmd/main"]