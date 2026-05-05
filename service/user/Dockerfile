FROM golang:1.25.3-alpine AS builder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache tzdata

# 在项目根目录构建（最稳）
WORKDIR /build

# 拷贝整个项目（支持 go.work）
COPY . .

# 下载依赖
RUN go mod download

# 用包路径构建（关键）
RUN go build -ldflags="-s -w" -o /app/user ./service/user


FROM scratch

# 证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 时区
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app

# 二进制
COPY --from=builder /app/user /app/user

CMD ["./user", "-f", "etc/user.yaml"]