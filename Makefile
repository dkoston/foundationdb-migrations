build:
	cwd=$(pwd 2>&1) && cd cmd/fdbm && go get . && go build -ldflags=-s . && cd ${cwd}

install: build
	cp cmd/fdbm/fdbm ${GOPATH}/bin/

test:
	cwd=$(pwd 2>&1) && cd lib/fdbm && go get . && go test ./... && cd ${cwd}

.PHONY: clean test all
clean:
	rm -f cmd/fdbm/fdbm