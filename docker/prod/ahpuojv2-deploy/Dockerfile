FROM ubuntu:18.04

COPY ./assets/sources.list /etc/apt/sources.list

RUN set -ex \
	&& export DEBIAN_FRONTEND=noninteractive \
	&& chmod -R 777 /tmp \
	&& apt-get clean \
	&& mv /var/lib/apt/lists /tmp \
	&& mkdir -p /var/lib/apt/lists/partial \
	&& apt-get clean \
	&& apt-get update \
	&& apt-get install -y tzdata

RUN set -ex \
	&& apt-get install -y \
	make flex g++ libmysqlclient-dev libmysql++-dev \
	nginx \
	openjdk-8-jdk \
	nano \
	python3.7

RUN set -ex \
	&& apt-get install -y \
	mysql-client-5.7 \
	busybox \
	wget \
	supervisor

RUN set -ex \
	&& apt-get install -y git
# RUN set -ex \
# 	&& rm -rf /var/lib/apt/lists/*

RUN set -ex \
	&& ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& dpkg-reconfigure -f noninteractive tzdata

WORKDIR /home/judge/


RUN set -ex \
	&& wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
	&&  tar -xvf go1.13.1.linux-amd64.tar.gz \
	&& mv go /usr/lib/go \
	&& rm go1.13.1.linux-amd64.tar.gz

COPY ./assets/default /etc/nginx/sites-enabled/default
COPY ./assets/nginx.conf /etc/nginx/nginx.conf
COPY ./assets/supervisor_ahpuoj.conf /etc/supervisor/conf.d/
COPY ./assets/docker-entrypoint.sh /usr/local/bin/
COPY ./assets/install /home/judge/install

ENV GOROOT /usr/lib/go
ENV GOPATH /home/judge/go

RUN set -ex \
	&& chmod +x /usr/local/bin/docker-entrypoint.sh \
	&& ln -s /usr/local/bin/docker-entrypoint.sh  /docker-entrypoint.sh \
	&& mkdir /home/judge/go \
	&& mkdir /home/judge/go/bin \
	&& mkdir /home/judge/go/src \
	&& mkdir /home/judge/go/pkg \
