FROM ubuntu:18.04

COPY sources.list /etc/apt/ustc-sources.list

RUN set -ex \
	&& export DEBIAN_FRONTEND=noninteractive \
	&& apt-get update \
	&& apt-get install -y tzdata
	
RUN set -ex \
	&& apt-get install -y \
		make flex g++ libmysqlclient-dev libmysql++-dev \
		nginx \
		openjdk-8-jdk \
		mysql-client-5.7 \
		nano \
		python3.7 \
		golang-1.10 \
		busybox

RUN set -ex \
	&& apt-get install -y git
# RUN set -ex \
# 	&& rm -rf /var/lib/apt/lists/*

RUN set -ex \
	&& ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& dpkg-reconfigure -f noninteractive tzdata
	
WORKDIR /home/judge/

COPY ./default /etc/nginx/sites-enabled/default

COPY ./docker-entrypoint.sh /usr/local/bin/

ENV GOROOT /usr/lib/go-1.10
ENV GOPATH /home/judge/go
ENV PATH $PATH:$GOPATH/bin

RUN set -ex \
	&& chmod +x /usr/local/bin/docker-entrypoint.sh \
	&& ln -s /usr/local/bin/docker-entrypoint.sh  /docker-entrypoint.sh \
	&& ln -s /usr/lib/go-1.10/bin/go  /usr/bin/go \
	&& mkdir /home/judge/go \
	&& mkdir /home/judge/go/bin \
	&& mkdir /home/judge/go/src \
	&& mkdir /home/judge/go/pkg \
	&& go get github.com/silenceper/gowatch \
