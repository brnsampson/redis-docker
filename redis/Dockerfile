FROM golang:stretch as builder

RUN mkdir /build

COPY src/ /go/src/

WORKDIR /go

RUN apt-get update; \
    apt-get install -y git; \
    rm -rf /var/lib/apt/lists/*

RUN go get github.com/go-redis/redis \
    && go get github.com/hashicorp/consul/api \
    && go build -i -o manage manage \
    && chmod +x manage


From redis:4.0.11

RUN apt-get update; \
    apt-get install -y curl unzip; \
    rm -rf /var/lib/apt/lists/*; \
    mkdir app;

WORKDIR /app

# Install consul agent
RUN export CONSUL_VERSION=1.2.1 \
    && export CONSUL_ARCH=linux_amd64 \
    && curl --retry 7 --fail -vo /tmp/consul_sha256sums "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS" \
    && export CONSUL_CHECKSUM=$( cat /tmp/consul_sha256sums | grep ${CONSUL_ARCH} | awk '{print $1}' ) \
    && curl --retry 7 --fail -vo /tmp/consul.zip "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_${CONSUL_ARCH}.zip" \
    && echo "${CONSUL_CHECKSUM}  /tmp/consul.zip" | sha256sum -c \
    && unzip /tmp/consul -d /usr/local/bin \
    && rm /tmp/consul.zip \
    && rm /tmp/consul_sha256sums \
    && mkdir /config

# Install ContainerPilot
ENV CONTAINERPILOT_VERSION 3.8.0
RUN curl -Lso /tmp/containerpilot.sha1.txt \
         "https://github.com/joyent/containerpilot/releases/download/${CONTAINERPILOT_VERSION}/containerpilot-${CONTAINERPILOT_VERSION}.sha1.txt" \
    && export CP_SHA1=$( cat /tmp/containerpilot.sha1.txt | grep containerpilot-${CONTAINERPILOT_VERSION}.tar.gz | awk '{print $1}' ) \
    && curl -Lso /tmp/containerpilot.tar.gz \
         "https://github.com/joyent/containerpilot/releases/download/${CONTAINERPILOT_VERSION}/containerpilot-${CONTAINERPILOT_VERSION}.tar.gz" \
    && echo "${CP_SHA1}  /tmp/containerpilot.tar.gz" | sha1sum -c \
    && tar zxf /tmp/containerpilot.tar.gz -C /bin \
    && rm /tmp/containerpilot.tar.gz \
    && rm /tmp/containerpilot.sha1.txt

# COPY ContainerPilot configuration
ENV CONTAINERPILOT_PATH=/etc/containerpilot.json5
COPY containerpilot.json5 ${CONTAINERPILOT_PATH}
ENV CONTAINERPILOT=${CONTAINERPILOT_PATH}

# Copy redis configuration
# COPY consul.hcl /etc/consul/
COPY --from=builder /go/manage /app/
COPY redis.conf /usr/local/etc/redis/redis.conf

ENV REDIS_ADDR 127.0.0.1:6379
ENV REDIS_PASS P@ssword!

VOLUME ["/data"]

ENTRYPOINT []

CMD ["/bin/containerpilot"]
