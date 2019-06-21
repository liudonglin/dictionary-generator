FROM node:8-slim as build-ui
WORKDIR /go/src/github.com/code-server/ui
COPY ./ui .
RUN npm config set registry https://registry.npm.taobao.org && npm install && npm run build

FROM golang:latest as builder-server
ENV GO15VENDOREXPERIMENT=1
ENV GO111MODULE=on
WORKDIR /go/src/github.com/code-server
COPY ./server .
RUN go build -mod=vendor -o server

FROM alpine as final
RUN apk --update upgrade && apk add sqlite
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder-server /go/src/github.com/code-server/server /app/server
COPY --from=build-ui /go/src/github.com/code-server/server/ui-dist /app/ui-dist

EXPOSE 8080

ENTRYPOINT ["/app/server"]

# 构建镜像
# docker build --rm=true -t xiayangqiushi/code-server:1.0 .

# 运行镜像
# docker run -d -it -p 8080:8080 xiayangqiushi/code-server:1.0