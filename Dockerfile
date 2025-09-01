# 多阶段构建 - 构造阶段
FROM golang:1.23.10-alpine3.22 AS builder

# 设置工作目录
WORKDIR /app

# 安装构建工具链
RUN apk add --no-cache --update \
    git \
    gcc \
    musl-dev \
    make \
    && rm -rf /var/cache/apk/*

# 下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目文件
COPY . .

# 编译应用（关键修复）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o /starmall .

# 多阶段构建 - 运行阶段
FROM alpine:3.22.0

# 系统配置
RUN apk add --no-cache ca-certificates tzdata curl \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

RUN addgroup -S appuser
# 创建应用用户
RUN adduser -D -h /app -s /sbin/nologin -G appuser appuser

# 创建应用目录结构
WORKDIR /app
RUN mkdir -p {conf,logs,storage} \
    && chown -R appuser:appuser /app

RUN mkdir -p /var/log/star-Mall
RUN touch /var/log/star-Mall/conf.log


# 配置文件复制
COPY --from=builder /starmall /app/starmall

# 权限设置
RUN chmod +x /app/starmall

# 启动应用
USER appuser
EXPOSE 8080
CMD ["/app/starmall"]