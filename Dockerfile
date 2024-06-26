FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod vendor

COPY . .

RUN go build -o main .

COPY mongodb_init.js /docker-entrypoint-initdb.d/

# RUN make download-file

RUN go run cli/convert.go

EXPOSE 8080
CMD ["./main"]
