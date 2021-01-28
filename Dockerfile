FROM golang:1.14-buster AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GO111MODULE=on
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

# 设置 gobuilder 中的工作路径，因为有依赖库的需求，所以需要手动配置需要在此拉去的 git 项目
# 创建编译的 root 目录，且拉去对应的支持库的代码
#RUN mkdir -p /homalab \
#    && mkdir -p /app \
#    && cd /homalab \
# 切换工作目录
WORKDIR /homalab/buildspace
# 添加到对应 WORKDIR 中，还原依赖
#ADD ./go.mod .
#ADD ./go.sum .
#RUN go mod download
# 把需要编译的代码一堆东西 COPY 到当前目录
COPY . .
# 执行编译，-o 指定保存位置和程序编译名称
RUN go build -ldflags="-s -w" -o /app/rssproxy

FROM alpine

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
# 主程序
COPY --from=builder /app/rssproxy /app/rssproxy
# 配置文件
COPY --from=builder /homalab/buildspace/config.yaml.sample /app/config.yaml
RUN chmod -R 777 /app
EXPOSE 1200

ENTRYPOINT ["/app/rssproxy"]