FROM golang:1.21.0
WORKDIR /web
COPY . .
RUN go build -o app
EXPOSE 8080
CMD ["./app"]