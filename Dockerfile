# 使用官方的 Go 语言镜像作为基础镜像
FROM bitnami/golang:1.24.1-debian-12-r1 AS builder

# 设置工作目录为 /app
# 所有后续操作都会在这个目录下进行
WORKDIR /opt/app

# 将当前项目目录的所有文件拷贝到容器的 /app 目录中
COPY . .

RUN sh apt-change-sources.sh

RUN apt update

RUN apt install libtagc0-dev -y

# 设置 Go 模块代理为 https://goproxy.cn（在中国加速模块下载），并下载项目的依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w CGO_ENABLED=1 && go mod download

# 编译 Go 项目，生成可执行文件 
RUN make build

# 使用一个更小的基础镜像（Alpine）来运行应用程序
# Alpine 是一个极简的 Linux 发行版，适合部署阶段
FROM debian:12-slim

# 设置工作目录为 /data/senanomusic/
RUN mkdir -p /data/senanomusic/cache && chmod -R 755 /data/senanomusic
WORKDIR /data/senanomusic

COPY --from=builder /usr/lib/x86_64-linux-gnu/libtag* /usr/lib/x86_64-linux-gnu/

# 从编译阶段的镜像中拷贝编译后的二进制文件到运行镜像中
COPY --from=builder /opt/app/bin/senanomusic /usr/local/bin/senanomusic
RUN chmod +x /usr/local/bin/senanomusic
COPY --from=builder /opt/app/bin/.env /data/senanomusic/.env

# 暴露容器的 8080 端口，用于外部访问
EXPOSE 8000

# 设置容器启动时运行的命令
# 这里是运行编译好的可执行文件 /usr/local/bin/senanomusic
CMD ["/usr/local/bin/senanomusic", "-env=/data/senanomusic/.env"]