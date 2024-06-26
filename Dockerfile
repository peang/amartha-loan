FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod vendor

COPY . .

RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
