#!/bin/bash

docker rmi myredis:6380.0
docker build -t myredis:6380.0 .
docker run --name redis6379 -d -v ./share6379:/myRedis:Z -p 6379:6379 myredis:6379.0
docker run --name redis6380 -d -v ./share6380:/myRedis:Z -p 6380:6380 myredis:6380.0
docker kill redis6379 && docker rm redis6379
docker kill redis6380 && docker rm redis6380
docker exec -it redis6379 /bin/bash
docker exec -it redis6380 /bin/bash
docker exec -it redis6379 redis-cli
docker exec -it redis6380 redis-cli
INFO replication
docker exec redis6379 redis-server ./conf/redis6379.conf
docker exec redis6380 redis-server ./conf/redis6380.conf