FROM ubuntu:latest

# install golang
RUN apt-get update && apt-get install -y golang && apt-get install -y go-dep

# install git
RUN apt-get update && apt-get install -y git

# install redis
RUN apt-get update && apt-get install -y redis

# install rabbitmq-server
RUN apt-get update && apt-get install -y rabbitmq-server

# install tools
RUN apt-get update && apt-get install -y vim curl net-tools

# install tools
RUN mkdir -p /root/go/bin

# configure bind adresses
RUN sed -i 's/bind 127.0.0.1 ::1/bind 127.0.0.1/' /etc/redis/redis.conf 
RUN sed -i  '/NODE_IP_ADDRESS/ s/^#NODE_IP_ADDRESS/NODE_IP_ADDRESS/g' /etc/rabbitmq/rabbitmq-env.conf

# enable management plugin on rabitmq
RUN rabbitmq-plugins enable rabbitmq_management;
