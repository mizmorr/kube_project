FROM golang:1.21
WORKDIR /app
COPY go.mod ./
COPY ./static/* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /listen
EXPOSE 8080
CMD ["/listen"]
