FROM golang:1.12.7 as builder

COPY . /code

WORKDIR /code

RUN unset GOPATH && \
    go test -v ./... && \
    go install ./...

FROM golang:1.12.7

RUN mkdir -p /opt/resource

COPY --from=builder /root/go/bin/* /opt/resource/
