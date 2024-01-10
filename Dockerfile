FROM golang:1.21-alpine as builder

RUN apk add --no-cache --update git make
RUN mkdir /build
WORKDIR /build
RUN git clone https://github.com/sebidude/kubeinfo.git
WORKDIR /build/kubeinfo
RUN make build-linux

FROM scratch

COPY --from=builder /build/kubeinfo/build/linux/kubeinfo /usr/bin/kubeinfo
ENTRYPOINT ["/usr/bin/kubeinfo"]
