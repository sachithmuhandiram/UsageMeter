FROM golang:alpine

RUN apk update

RUN apk add git

COPY . /go/src

WORKDIR /go/src

ENV CGO_ENABLED=0 GO111MODULE=off

RUN go build -o main .

EXPOSE 7272

RUN chmod 755 main

CMD [ "./main" ]

#RUN go test -v