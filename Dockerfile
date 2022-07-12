FROM golang:1.17

ENV GO111MODULE="on"

ENV GOPROXY="https://goproxy.cn"

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  apt-get update

RUN  go build -o main

EXPOSE 1926

CMD ["./main"]