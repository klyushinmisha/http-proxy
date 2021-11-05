# http-proxy

## Configurable HTTP proxy with distributed deadlock detection

## Basic usage

```
$ docker run --rm -v $(pwd)/example/config.yaml:/etc/http-proxy/config.yaml --network host registry.gitlab.com/klyushinmisha/http-proxy:0.1.0
2021/11/05 20:05:46 

 /$$   /$$ /$$$$$$$$ /$$$$$$$$ /$$$$$$$        /$$$$$$$                                        
| $$  | $$|__  $$__/|__  $$__/| $$__  $$      | $$__  $$                                       
| $$  | $$   | $$      | $$   | $$  \ $$      | $$  \ $$ /$$$$$$   /$$$$$$  /$$   /$$ /$$   /$$
| $$$$$$$$   | $$      | $$   | $$$$$$$/      | $$$$$$$//$$__  $$ /$$__  $$|  $$ /$$/| $$  | $$
| $$__  $$   | $$      | $$   | $$____/       | $$____/| $$  \__/| $$  \ $$ \  $$$$/ | $$  | $$
| $$  | $$   | $$      | $$   | $$            | $$     | $$      | $$  | $$  >$$  $$ | $$  | $$
| $$  | $$   | $$      | $$   | $$            | $$     | $$      |  $$$$$$/ /$$/\  $$|  $$$$$$$
|__/  |__/   |__/      |__/   |__/            |__/     |__/       \______/ |__/  \__/ \____  $$
                                                                                      /$$  | $$
                                                                                     |  $$$$$$/
                                                                                      \______/ 

2021/11/05 20:05:46 host: 0.0.0.0
port: 8080
hosts: [localhost:8081 localhost:8082]
balancer-type: round-robin
max-header: 10240 bytes
slab-size: 4096 bytes
max-concurrent-requests: 4
```

## Load balancing

In case of round-robin balancer (i.e. two different hosts return PING or PONG)
```
$ http localhost:8080
HTTP/1.1 200 OK
Content-Length: 4
Content-Type: text/plain; charset=utf-8
Date: Fri, 05 Nov 2021 20:10:19 GMT
Server: http-proxy/0.1.0
X-Proxy-Id: eb332249-6177-41bd-a1b1-655b08b23096

PING


$ http localhost:8080
HTTP/1.1 200 OK
Content-Length: 4
Content-Type: text/plain; charset=utf-8
Date: Fri, 05 Nov 2021 20:10:20 GMT
Server: http-proxy/0.1.0
X-Proxy-Id: eb332249-6177-41bd-a1b1-655b08b23096

PONG


$ http localhost:8080
HTTP/1.1 200 OK
Content-Length: 4
Content-Type: text/plain; charset=utf-8
Date: Fri, 05 Nov 2021 20:10:21 GMT
Server: http-proxy/0.1.0
X-Proxy-Id: eb332249-6177-41bd-a1b1-655b08b23096

PING
```

## Errors handling

If servers aren't reachable

```
2021/11/05 20:07:29 Options "http://localhost:8081": dial tcp 127.0.0.1:8081: connect: connection refused
```

If circuit detected

```
2021/11/05 20:08:24 distributed deadlock detected
```
