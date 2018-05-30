FROM golang:1.10-stretch

WORKDIR /go/src/github.com/dkoston/foundationdb-migrations
RUN     mkdir -p /go/src/github.com/dkoston/foundationdb-migrations
COPY    Makefile /go/src/github.com/dkoston/foundationdb-migrations/Makefile
COPY    db  /go/src/github.com/dkoston/foundationdb-migrations/db
COPY    cmd /go/src/github.com/dkoston/foundationdb-migrations/cmd
COPY    lib /go/src/github.com/dkoston/foundationdb-migrations/lib
#RUN     apt-get update \
#            && apt-get install -y apt-transport-https dirmngr m4 \
#            && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF \
#            && echo "deb https://download.mono-project.com/repo/debian stable-stretch main" | tee /etc/apt/sources.list.d/mono-official-stable.list \
#            && apt-get update \
#            && apt-get install -y mono-devel \
#            && wget https://www.foundationdb.org/downloads/5.1.7/ubuntu/installers/foundationdb-clients_5.1.7-1_amd64.deb \
#            && dpkg -i foundationdb-clients_5.1.7-1_amd64.deb \
#            && cd /root \
#            && git clone https://github.com/apple/foundationdb.git \
#            && cd foundationdb/bindings/go \
#            && ./fdb-go-install.sh install \
#            && cd /go/src/github.com/dkoston/foundationdb-migrations \
#            && go get ./... \
#            && make install