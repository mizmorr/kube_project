FROM golang:1.21
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./main.go
EXPOSE 8080
CMD ["/main"]
