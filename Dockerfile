FROM golang:1.23.0

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o kalkan ./cmd

EXPOSE 8080

CMD ["./kalkan"]