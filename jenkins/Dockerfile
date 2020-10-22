# 基础镜像
FROM golang:alpine

# 工作目录
WORKDIR /

# 在mac上开发的，需要用到go的交叉编译
# RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags="-s -w" -o quickpass

# 复制编译后的main文件
COPY ./main /main

# 复制配置文件目录
COPY ./conf /conf

# 执行命令
CMD ["./main"]