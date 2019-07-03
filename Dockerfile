FROM golang:1.12

COPY . /go/src/github.com/t1labs/aps

WORKDIR /go/src/github.com/t1labs/aps

RUN go install github.com/t1labs/aps/cmd/aps

CMD ["aps"]