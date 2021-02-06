FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/lib/pq

ADD . /go/src/app
WORKDIR /go/src/app

EXPOSE 80

CMD ["go","run","main.go"]