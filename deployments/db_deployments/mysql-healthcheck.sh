#!/bin/bash

# MySQL 连接信息
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASSWORD="o?dv(%qn)Uf8*e^3AQTPB4k6L1N975xG"
MYSQL_DATABASE="web"

# 检查 MySQL 连接
mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" -e "SELECT 1;" > /dev/null 2>&1

# 检查 MySQL 连接结果并输出相应状态码
if [ $? -eq 0 ]; then
    echo "MySQL is healthy"
    exit 0
else
    echo "MySQL is not healthy"
    exit 1
fi