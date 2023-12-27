FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8081

CMD ["go", "run", "main.go"]
