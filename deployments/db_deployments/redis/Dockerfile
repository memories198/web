FROM redis:7.2.4
LABEL authors="memories"

WORKDIR /myRedis

USER redis

ENTRYPOINT ["export", "TZ=Asia/Shanghai"]
ENTRYPOINT ["redis-server", "conf/redis6380.conf"]