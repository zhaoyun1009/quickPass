FROM golang:1.14 as build-stage
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOARM=6
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download && go mod verify && go build -ldflags="-s -w" -o quickpass

FROM alpine:3.11
WORKDIR /app
RUN mkdir conf
COPY --from=build-stage /app/quickpass .
COPY --from=build-stage /app/conf/* /app/conf/
RUN rm /app/conf/app.ini && mv /app/conf/app_iopay666.ini /app/conf/app.ini
RUN echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories \
    && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
EXPOSE 8000
CMD ["./quickpass"]
