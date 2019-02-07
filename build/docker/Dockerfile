FROM ubuntu:latest

# install redis
RUN apt-get update && apt-get install -y redis

# install rabbitmq-server
RUN apt-get update && apt-get install -y rabbitmq-server

# install tools
RUN apt-get update && apt-get install -y vim curl net-tools

# configure bind adresses
RUN sed -i 's/bind 127.0.0.1 ::1/bind 127.0.0.1/' /etc/redis/redis.conf 
RUN sed -i  '/NODE_IP_ADDRESS/ s/^#NODE_IP_ADDRESS/NODE_IP_ADDRESS/g' /etc/rabbitmq/rabbitmq-env.conf

# enable management plugin on rabitmq
RUN rabbitmq-plugins enable rabbitmq_management;

# Copy the binary
COPY uniris-miner ./root

# Configure the app
ENV UNIRIS_PUBLIC_KEY=3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
ENV UNIRIS_PRIVATE_KEY=307702010104201432f00062c0229d19e24c070a713d900da4883788ce3f8bd3fede4c10a36a79a00a06082a8648ce3d030107a1440342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
ENV UNIRIS_NETWORK_TYPE=private
ENV UNIRIS_NETWORK_INTERFACE=eth0
ENV UNIRIS_DISCOVERY_SEEDS=127.0.0.1:5000:3059301306072a8648ce3d020106082a8648ce3d0301070342000459c8b568df66798d7f876d94fb0afc516502893d996610632c40f70b830aebf39e0cbee311af4450ec56859d2b8f59ec09a44c7e303d030899aee551de61af2e
ENV UNIRIS_DISCOVERY_DB_TYPE=redis
ENV UNIRIS_DISCOVERY_DB_HOST=localhost
ENV UNIRIS_DISCOVERY_DB_PORT=6379
ENV UNIRIS_DISCOVERY_DB_PWD=
ENV UNIRIS_BUS_TYPE=amqp
ENV UNIRIS_BUS_HOST=localhost
ENV UNIRIS_BUS_PORT=5672
ENV UNIRIS_BUS_USER=guest
ENV UNIRIS_BUS_PWD=guest
ENV UNIRIS_INT_GRPC_PORT=3009
ENV UNIRIS_EXT_GRPC_PORT=5000

# Expose Public services
EXPOSE 4000
EXPOSE $UNIRIS_EXT_GRPC_PORT

#Expose Rabitmq management console
EXPOSE 15672

WORKDIR /root

#Run the app
CMD /etc/init.d/redis-server start ; \
    /etc/init.d/rabbitmq-server start ; \
    rabbitmqctl add_user uniris uniris; \
    rabbitmqctl set_user_tags uniris administrator; \
    rabbitmqctl set_permissions -p / uniris ".*" ".*" ".*"; \
    ./uniris-miner