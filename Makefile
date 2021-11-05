.PHONY: all build clean

all: clean build

build: http-proxy

clean:
	rm -f http-proxy

http-proxy:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-proxy cmd/http-server/main.go
