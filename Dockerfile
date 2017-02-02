FROM golang:1.6

RUN mkdir -p /go/src/imagesServer
WORKDIR /go/src/imagesServer

COPY . /go/src/imagesServer
RUN mkdir -p static/multipleFiles
RUN mkdir -p static/files
RUN go get -d -v
RUN go install -v
RUN go build

CMD ["./imagesServer"]
