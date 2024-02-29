FROM ubuntu:latest
LABEL authors="memories"

WORKDIR /root
COPY ./cpu.sh .
RUN chmod u+x cpu.sh
ENTRYPOINT ["/root/cpu.sh"]

#docker build -t cpu -f Dockerfile --no-cache .
#docker run --name cpu --rm -d cpu
#