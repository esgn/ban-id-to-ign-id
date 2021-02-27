FROM golang:alpine

ENV GO111MODULE=auto

RUN apk update && apk add git

RUN go version

ADD . /go/src/app
WORKDIR /go/src/app

RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

ENV CGO_ENABLED=0
RUN go test -v
ENV CGO_ENABLED=1

EXPOSE 8080

CMD ["go","run","api.go"]
