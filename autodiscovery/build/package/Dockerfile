FROM uniris:latest

# Install the dependencies
ENV GOPATH=/root/go
WORKDIR $GOPATH/src/github.com/uniris/uniris-core/autodiscovery
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Copy the code
COPY . ./ 

# Configure the app
ENV UNIRIS_VERSION=1.0
ENV UNIRIS_PUBLICKEY="MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhNngAzz3Ws1cz7veHnzeCwi9cuJhmpcHoi8+q3y75UL1vG5hyfnEMTjDCeeRHiobASUGvawGlcxexJzlZPIVYg=="
ENV UNIRIS_NETWORK=private
ENV UNIRIS_NETWORK_INTERFACE=lo0
ENV UNIRIS_DISCOVERY_PORT=3545
ENV UNIRIS_DISCOVERY_SEEDS=127.0.0.1:3545:$UNIRIS_PUBLICKEY
ENV UNIRIS_DISCOVERY_REDIS_HOST=localhost
ENV UNIRIS_DISCOVERY_REDIS_PORT=6379
ENV UNIRIS_DISCOVERY_REDIS_PWD=
ENV UNIRIS_DISCOVERY_AMQP_HOST=localhost
ENV UNIRIS_DISCOVERY_AMQP_PORT=5672
ENV UNIRIS_DISCOVERY_AMQP_USER=guest
ENV UNIRIS_DISCOVERY_AMQP_PWD=guest

EXPOSE $UNIRIS_DISCOVERY_PORT

# Run the app
CMD sed -i 's/bind 127.0.0.1 ::1/bind 127.0.0.1/' /etc/redis/redis.conf ; /etc/init.d/redis-server start ; /etc/init.d/rabbitmq-server start ; go run cmd/discovery/main.go