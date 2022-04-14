FROM golang:1.18.1-alpine as builder
LABEL author="Thomas Bellembois"
ARG BuildID
ENV BuildID=${BuildID}

# Install GCC and git.
RUN apk add build-base git

# ref. go.mod gochimitheque-wasm
RUN mkdir -p /home/thbellem/workspace \
    && ln -s /go /home/thbellem/workspace/workspace_go

# Installing dependencies.
RUN go install github.com/Joker/jade/cmd/jade@master

#
# Sources.
#

# Getting wasm module sources.
WORKDIR /go/src/github.com/tbellembois/
# sudo mount --bind ~/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-wasm ./bind-gochimitheque-wasm
# sudo mount --bind ~/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils ./bind-gochimitheque-utils
COPY ./bind-gochimitheque-wasm ./gochimitheque-wasm
COPY ./bind-gochimitheque-utils ./gochimitheque-utils

# Copying Chimithèque sources.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
COPY . .

#
# Build.
#

# Building wasm module.
WORKDIR /go/src/github.com/tbellembois/gochimitheque-wasm
RUN GOOS=js GOARCH=wasm go get -v -d ./... \
    && GOOS=js GOARCH=wasm go build -o wasm .

# Copying and compress WASM module into sources.
RUN cp /go/src/github.com/tbellembois/gochimitheque-wasm/wasm /go/src/github.com/tbellembois/gochimitheque/wasm/ \
    && gzip -9 -v -c /go/src/github.com/tbellembois/gochimitheque/wasm/wasm > /go/src/github.com/tbellembois/gochimitheque/wasm/wasm.gz \
    && rm /go/src/github.com/tbellembois/gochimitheque/wasm/wasm

# Installing Chimithèque dependencies.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/

# Generating code.
RUN go generate

# Building Chimithèque.
# docker build --build-arg BuildID=2.0.7 -t tbellembois/gochimitheque:2.0.7 .
RUN if [ -z $BuildID ]; then BuildID=$(date "+%Y%m%d"); fi; echo "BuildID=$BuildID"; go build -ldflags "-X main.BuildID=$BuildID"

#
# Install.
#

FROM golang:1.18beta2-alpine

RUN apk add bash && rm -Rf /var/cache/apk

# Ensure www-data user exists.
RUN set -x ; \
  addgroup -g 82 -S www-data ; \
  adduser -u 82 -D -S -G www-data www-data && exit 0 ; exit 1 \
  && mkdir /data \
  && chown www-data /data \
  && chmod 700 /data \
  && mkdir /var/www-data \
  && chown www-data /var/www-data

COPY --from=builder /go/src/github.com/tbellembois/gochimitheque/gochimitheque /var/www-data/
RUN chown www-data /var/www-data/gochimitheque \
    && chmod +x /var/www-data/gochimitheque

#
# Final work.
#

# Copying entrypoint.
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

# Adding ENS-Lyon CA certificates.
ADD docker/terena.crt /usr/local/share/ca-certificates/terena.crt
RUN chmod 644 /usr/local/share/ca-certificates/terena.crt && update-ca-certificates

# Container configuration.
USER www-data
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME [ "/data" ]
EXPOSE 8081