FROM golang:latest
WORKDIR /web
COPY . .
RUN go build -o app
EXPOSE 8080
CMD ["./app"]