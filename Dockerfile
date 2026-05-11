FROM golang:1.22.4-alpine
WORKDIR /logger

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o logger .
RUN chmod +x logger
CMD ["./logger"]