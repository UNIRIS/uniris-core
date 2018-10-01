FROM golang:1.8

# Download and install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install the dependencies
WORKDIR $GOPATH/src/github.com/uniris/uniris-core/autodiscovery
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Copy the code
COPY . ./ 

# Configure the app
ENV UNIRIS_VERSION=1.0
ENV UNIRIS_PUBLICKEY="MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEIh7hIRcsQsZJUFpWWNzrLwq53exmXKi7duDuts8L4xasepdT6GxTQOn9C/0dzwcah7bqiRlcpNpYrKJyBetTqQ=="
ENV UNIRIS_NETWORK=private
ENV UNIRIS_NETWORK_INTERFACE=lo0
ENV UNIRIS_DISCOVERY_PORT=3545
ENV UNIRIS_DISCOVERY_P2PFACTOR=1
ENV UNIRIS_DISCOVERY_SEEDS=127.0.0.1:3545:MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhNngAzz3Ws1cz7veHnzeCwi9cuJhmpcHoi8+q3y75UL1vG5hyfnEMTjDCeeRHiobASUGvawGlcxexJzlZPIVYg==
ENV UNIRIS_DISCOVERY_REDIS_HOST=localhost
ENV UNIRIS_DISCOVERY_REDIS_PORT=6379
ENV UNIRIS_DISCOVERY_REDIS_PWD=
ENV UNIRIS_DISCOVERY_AMQP_HOST=localhost
ENV UNIRIS_DISCOVERY_AMQP_PORT=5672
ENV UNIRIS_DISCOVERY_AMQP_USER=guest
ENV UNIRIS_DISCOVERY_AMQP_PWD=guest

EXPOSE $UNIRIS_DISCOVERY_PORT

# Run the app
CMD go run cmd/discovery/main.go