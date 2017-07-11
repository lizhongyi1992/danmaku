.PHONY: all clean

all:
	go build -ldflags "-X main.VERSION=1.0.0 -X 'main.GITHASH=`git rev-parse --short HEAD`' -X 'main.BUILT=`date +"%Y/%m/%d %H:%M:%S"`' -X 'main.GOVERSION=`go version | cut -d" " -f 3`'"
	./counter -genconf > config.yaml

clean:
	go clean
