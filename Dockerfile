FROM golang:1.15-stretch as builder
WORKDIR /dongtzu-server
ADD . /dongtzu-server
RUN go build

FROM centos:centos7
WORKDIR /app
COPY --from=builder /dongtzu-server/dongtzu .
COPY ./conf.d/env.template.yaml ./conf/env.yaml
CMD ["./dongtzu","server","-f","./conf/env.yaml"]