FROM ubuntu
LABEL authors="memories"

COPY ./web /web/
WORKDIR /web
RUN chmod u+x web

# 添加一个普通用户, -r 选项是 groupadd 命令的一个参数，它表示创建一个系统用户组（system group）。系统用户组（system group）是一种特殊类型的用户组，通常用于系统进程或服务运行时的权限控制。与普通用户组不同，系统用户组的 GID（Group Identifier，用户组标识符）通常是较小的数字，并且不会被分配给普通用户。
RUN groupadd -r webUser && useradd -r -g webUser webUser
RUN chown -R webUser:webUser /web
USER webUser
RUN ln -s /share/config /web/config
RUN ln -s /share/logs  /web/logs

EXPOSE 80

CMD ["./web"]