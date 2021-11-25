FROM golang:alpine

WORKDIR /server
COPY . .

RUN go mod download
RUN go mod verify

RUN go build -o server cmd/main.go

ENTRYPOINT ["./server"]