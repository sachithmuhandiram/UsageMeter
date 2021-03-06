FROM golang:alpine

RUN apk update

RUN apk add git

RUN go get github.com/go-sql-driver/mysql

COPY . /go/src

WORKDIR /go/src

ENV CGO_ENABLED=0 GO111MODULE=off

#RUN go vet

RUN go build -o main .

EXPOSE 7171

RUN chmod 755 main

CMD [ "./main" ]

#RUN go test -v