# 基础镜像
FROM ubuntu:20.04

# 编译后的复制打包进镜像
COPY webook /app/webook

# 工作目录
WORKDIR /app

# CMD 是执行命令
# 启动（最佳）
ENTRYPOINT ["/app/webook"]
