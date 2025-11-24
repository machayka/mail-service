FROM golang:1.25.3-alpine
WORKDIR /app
COPY . .
RUN go build -o bin/app cmd/app/main.go
EXPOSE 420
CMD ["/app/bin/app"]
