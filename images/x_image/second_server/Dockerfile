FROM golang:1.21
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server_second ./main.go
EXPOSE 8080
CMD ["/server_second"]
