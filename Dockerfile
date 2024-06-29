FROM golang:1.21

WORKDIR /app

ARG PORT

COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod vendor

COPY . .

RUN go build -o main .

COPY mongodb_init.js /docker-entrypoint-initdb.d/

EXPOSE $PORT

CMD ["./main"]
