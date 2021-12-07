FROM golang:alpine

RUN apk add git
WORKDIR /server
COPY . .

ARG LOCAL_REPO
RUN if [ "$LOCAL_REPO" = "off" ] ; then echo "build submodule" ; else RUN git submodule update --init --recursive ; fi

RUN go mod download
RUN go mod verify

RUN go build -o server cmd/main.go

ENTRYPOINT ["./server"]