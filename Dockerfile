FROM golang:1.19.5-alpine3.17
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080

CMD ["/app/main"]
