FROM uniris:latest

# Install the dependencies
ENV GOPATH=/root/go
WORKDIR $GOPATH/src/github.com/uniris/uniris-core
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

#Expose autodiscovery port
EXPOSE $UNIRIS_DISCOVERY_PORT

#Expose API port
EXPOSE 8080

#Expose DataMining Port
EXPOSE 3547

#Expose Rabitmq management console
EXPOSE 15672

# Run the app
CMD /etc/init.d/redis-server start ; \
/etc/init.d/rabbitmq-server start ; \
rabbitmqctl add_user uniris uniris; \
rabbitmqctl set_user_tags uniris administrator; \
rabbitmqctl set_permissions -p / uniris ".*" ".*" ".*"; \
go run autodiscovery/cmd/discovery/main.go > /var/log/uniris/discovery.log 2>&1 ; \
go run api/cmd/api/main.go > /var/log/uniris/api.log 2>&1 ; \
go run datamining/cmd/datamining/main.go > /var/log/uniris/datamining.log 2>&1 ; \


