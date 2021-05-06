FROM golang:1.16-buster
LABEL author="Thomas Bellembois"
ARG GIT_COMMIT
ENV ENV_GIT_COMMIT=$GIT_COMMIT

#
# Build prepare.
#

# ref. go.mod gochimitheque-wasm
RUN mkdir -p /home/thbellem/workspace
RUN ln -s /go /home/thbellem/workspace/workspace_go

# Creating DB volume directory.
RUN mkdir /data && chown www-data /data

# Creating www directory.
RUN mkdir /var/www-data && chown www-data /var/www-data

# Installing Jade command.
RUN go get -v github.com/Joker/jade/cmd/jade@master

#
# Sources.
#

# Getting wasm module sources.
WORKDIR /go/src/github.com/tbellembois/
# sudo mount --bind ~/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-wasm ./bind-gochimitheque-wasm
COPY ./bind-gochimitheque-wasm ./gochimitheque-wasm

# Copying Chimithèque sources.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
COPY . .
COPY .git/ ./.git/

#
# Build.
#

# Building wasm module.
WORKDIR /go/src/github.com/tbellembois/gochimitheque-wasm
RUN GOOS=js GOARCH=wasm go get -v -d ./...
RUN GOOS=js GOARCH=wasm go build -o wasm .

# Copying and compress WASM module into sources.
RUN cp /go/src/github.com/tbellembois/gochimitheque-wasm/wasm /go/src/github.com/tbellembois/gochimitheque/wasm/
RUN gzip -9 -v -c /go/src/github.com/tbellembois/gochimitheque/wasm/wasm > /go/src/github.com/tbellembois/gochimitheque/wasm/wasm.gz
RUN rm /go/src/github.com/tbellembois/gochimitheque/wasm/wasm

# Installing Chimithèque dependencies.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
RUN go get -v -d ./...

# Generating code.
RUN go generate

# Building Chimithèque.
# docker build --build-arg GIT_COMMIT=2.0.7 -t tbellembois/gochimitheque:2.0.7 .
RUN if [ ! -z "$ENV_GIT_COMMIT" ]; then export GIT_COMMIT=$ENV_GIT_COMMIT; else export GIT_COMMIT=$(git rev-list -1 HEAD); fi; echo "version=$GIT_COMMIT" ;go build -ldflags "-X main.GitCommit=$GIT_COMMIT"

#
# Install.
#

# Installing Chimithèque.
RUN cp /go/src/github.com/tbellembois/gochimitheque/gochimitheque /var/www-data/ \
    && chown www-data /var/www-data/gochimitheque \
    && chmod +x /var/www-data/gochimitheque

#
# Final work.
#

# Cleaning up sources.
RUN rm -Rf /go/src/*

# Copying entrypoint.
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

# Container configuration.
USER www-data
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME ["/data"]
EXPOSE 8081