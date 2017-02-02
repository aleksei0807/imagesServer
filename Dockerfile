FROM golang:1.6

RUN mkdir -p /go/src/imagesServer
WORKDIR /go/src/imagesServer

# this will ideally be built by the ONBUILD below ;)

COPY . /go/src/imagesServer
RUN mkdir -p static/multipleFiles
RUN mkdir -p static/files
RUN go get -d -v
RUN go install -v
RUN go build

CMD ["./imagesServer"]
