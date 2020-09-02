FROM golang:latest
WORKDIR /root/micro-go-course/section10/user
COPY / /root/micro-go-course/section10/user
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
EXPOSE 10086
ENTRYPOINT ["./user"]