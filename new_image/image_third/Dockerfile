FROM golang:1.21
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /hello
EXPOSE 8080
CMD ["/hello"]
