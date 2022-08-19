FROM golang:1.17-alpine AS builder

# 为镜像设置必要的环境变量

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.io,direct \
    GOARCH=amd64

# 移动到工作目录

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# 将当前dockerfile所在目录下所有文件移动到当前工作目录
# 也就是代码复制到容器中 /build目录下
COPY . .

# 这一步编译了二进制文件
RUN go build -o bluebell .

# 重新依赖一个小镜像镜像
FROM scratch

# 把静态资源和配置文件拷贝过来
COPY ./config.yaml /
COPY ./templates /templates
COPY ./static /static

# 拷贝构建好的二进制执行文件
COPY --from=builder /build/bluebell /

# 声明端口
EXPOSE 8888

ENTRYPOINT [ "/bluebell" ]
