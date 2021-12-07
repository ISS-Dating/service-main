FROM golang:alpine

RUN apk add git
WORKDIR /server
COPY . .

RUN git submodule update --init --recursive

RUN go mod download
RUN go mod verify

RUN go build -o server cmd/main.go

ENTRYPOINT ["./server"]