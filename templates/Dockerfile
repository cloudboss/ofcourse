FROM golang:1.11.4 as builder

COPY . /code

WORKDIR /code

RUN unset GOPATH && \
    go test -v ./... && \
    go install ./...

FROM golang:1.11.4

RUN mkdir -p /opt/resource

COPY --from=builder /root/go/bin/* /opt/resource/
