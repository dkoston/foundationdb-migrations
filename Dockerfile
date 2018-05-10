FROM gcr.io/cryptowalk-price-alerts/golang:1.9

# Install goose

WORKDIR /go/src/github.com/cryptowalkio/goose
RUN     mkdir -p /go/src/github.com/cryptowalkio/goose
COPY    Makefile /go/src/github.com/cryptowalkio/goose/Makefile
COPY    db  /go/src/github.com/cryptowalkio/goose/db
COPY    cmd /go/src/github.com/cryptowalkio/goose/cmd
COPY    lib /go/src/github.com/cryptowalkio/goose/lib
COPY    test /go/src/github.com/cryptowalkio/goose/test
RUN     cd /go/src/github.com/cryptowalkio/goose \
    &&  chmod +x test/migration_test.bash \
    &&  make install


# Install psql for verifying migrations worked
RUN apt-get update \
    && apt-get install -y --no-install-recommends --force-yes postgresql-client \
    && echo 'postgresql-db:5432:database:test:abcd1234' > ~/.pgpass \
    && chmod 0600 ~/.pgpass

# Install node for testing node migrations

ENV NODE_VERSION 8.11.1

RUN set -ex \
  && curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz" \
  && tar -xJf "node-v$NODE_VERSION-linux-x64.tar.xz" -C /usr/local --strip-components=1 \
  && rm "node-v$NODE_VERSION-linux-x64.tar.xz" \
  && npm i -g npm@latest \
  && wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 \
  && chmod +x /usr/local/bin/dumb-init


RUN node -v
RUN npm -v
RUN npm config set progress false && npm config set loglevel warn