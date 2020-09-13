FROM golang:latest
WORKDIR /root
COPY / /root
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o register
EXPOSE 12312
ENTRYPOINT ./register -consul.addr=$consulAddr -service.addr=$serviceAddr