FROM scratch

COPY http-proxy /bin/http-proxy

ENTRYPOINT ["http-proxy"]
