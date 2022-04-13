FROM golang:1.16

WORKDIR /workdir/
ADD . /workdir/

RUN make go_get

CMD ["make", "go_test_integration"]
