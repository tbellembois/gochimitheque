# Version: 0.0.1
FROM golang:1.14-buster
LABEL author="Thomas Bellembois"

# copying sources
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
COPY . .

# installing GopherJS go1.12 dependency
RUN go get golang.org/dl/go1.12.16
RUN go1.12.16 download
ENV GOPHERJS_GOROOT=/root/sdk/go1.12.16

# installing dependencies
RUN go get -v ./...

# installing Chimithèque
RUN mkdir /var/www-data \
    && cp /go/bin/gochimitheque /var/www-data/ \
    && chown -R www-data /var/www-data \
    && chmod +x /var/www-data/gochimitheque

# cleanup sources
RUN rm -Rf /go/src/*

# copying entrypoint
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

# creating volume directory
RUN mkdir /data

USER www-data
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME ["/data"]
EXPOSE 8081
