#!/bin/bash

# Redis 连接信息
REDIS_HOST="127.0.0.1"
REDIS_PORT="6379"
REDIS_PASSWORD="o?dv(%qn)Uf8*e^3AQTPB4k6L1N975xG"

# 检查 Redis 连接
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" PING > /dev/null 2>&1

# 检查 Redis 连接结果并输出相应状态码
if [ $? -eq 0 ]; then
    echo "Redis is healthy"
    exit 0
else
    echo "Redis is not healthy"
    exit 1
fi