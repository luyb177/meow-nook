FROM golang:1.25.3-alpine AS builder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache tzdata

WORKDIR /build

COPY . .

RUN go mod download

#  关键修复
RUN go build -ldflags="-s -w" -o /app/gateway ./service/gateway


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/gateway /app/gateway

CMD ["./gateway", "-f", "etc/gateway.yaml"]