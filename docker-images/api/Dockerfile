FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

ADD . /go/src/app
WORKDIR /go/src/app

EXPOSE 8080

CMD ["go","run","api.go"]
