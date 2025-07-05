# 使用官方提供的dockerfile编译的基础镜像
FROM golang-alpine-custom:latest AS builder

# 设置工作目录为 /opt/app
# 所有后续操作都会在这个目录下进行
RUN mkdir -p /opt/app/ && chmod -R 755 /opt/app/
WORKDIR /opt/app/

# 将当前项目目录的所有文件拷贝到容器的 /opt/app 目录中
COPY . .

# 下载依赖并编译为可执行文件
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w CGO_ENABLED=1 && go mod download && make build

# 打包为可执行镜像
FROM alpine:3.21

# 安装taglib并创建工作目录
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk add --no-cache taglib; \
    mkdir -p /opt/app/ && chmod -R 755 /opt/app/

WORKDIR /opt/app/

# 从编译阶段的镜像中拷贝编译后的二进制文件到运行镜像中
COPY --from=builder /opt/app/target/bin/* /opt/app/

# 暴露容器的 8000 端口，用于外部访问
EXPOSE 8000

# 设置容器启动时运行的命令，运行 /opt/app/senanomusic
ENTRYPOINT ["./senanomusic", "-env=.env"]