FROM golang:1.21-bullseye as builder
LABEL author="Thomas Bellembois"
ARG BuildID
ENV BuildID=${BuildID}

# Install GCC and git.
# RUN apk add build-base git

# Install zeromq repository and library.
RUN echo 'deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_11/ /' | tee /etc/apt/sources.list.d/network:messaging:zeromq:release-stable.list
RUN curl -fsSL https://download.opensuse.org/repositories/network:messaging:zeromq:release-stable/Debian_11/Release.key | gpg --dearmor | tee /etc/apt/trusted.gpg.d/network_messaging_zeromq_release-stable.gpg > /dev/null
RUN apt -y update
RUN apt -y install libzmq3-dev

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

# Install Rust.
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="$PATH:/root/.cargo/bin"

# Get sources.
WORKDIR /go/src/rust
RUN git clone https://github.com/tbellembois/chimitheque_types.git
RUN git clone https://github.com/tbellembois/chimitheque_utils.git
RUN git clone https://github.com/tbellembois/chimitheque_utils_service.git

# Compile.
WORKDIR /go/src/rust/chimitheque_utils_service
RUN cargo build --release

#
# Install.
#

FROM golang:1.21-bullseye

# Install zeromq repository and library.
RUN echo 'deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_11/ /' | tee /etc/apt/sources.list.d/network:messaging:zeromq:release-stable.list
RUN curl -fsSL https://download.opensuse.org/repositories/network:messaging:zeromq:release-stable/Debian_11/Release.key | gpg --dearmor | tee /etc/apt/trusted.gpg.d/network_messaging_zeromq_release-stable.gpg > /dev/null
RUN apt -y update
RUN apt -y install libzmq3-dev

RUN rm -Rf /var/cache/apk

# Ensure www-data user exists.
RUN addgroup --gid 82 --system chimitheque \
  && adduser --uid 82 --system --ingroup chimitheque chimitheque \
  && mkdir /data \
  && chown chimitheque /data \
  && chmod 700 /data \
  && mkdir /var/www-data \
  && chown chimitheque /var/www-data \
  && chown chimitheque /var/log \
  && chmod 755 /var/log

COPY --from=builder /go/src/github.com/tbellembois/gochimitheque/gochimitheque /var/www-data/
RUN chown chimitheque /var/www-data/gochimitheque \
    && chmod +x /var/www-data/gochimitheque

COPY --from=builder /go/src/rust/chimitheque_utils_service/target/release/chimitheque_utils_service /var/www-data/
RUN chown chimitheque /var/www-data/chimitheque_utils_service \
    && chmod +x /var/www-data/chimitheque_utils_service

#
# Final work.
#

# Copying entrypoint.
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

# Adding CA certificates.
ADD docker/terena.crt /usr/local/share/ca-certificates/terena.crt
ADD docker/USERTrust_RSA_Certification_Authority.crt /usr/local/share/ca-certificates/USERTrust_RSA_Certification_Authority.crt

RUN chmod 644 /usr/local/share/ca-certificates/* && update-ca-certificates

# Container configuration.
USER chimitheque
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME [ "/data" ]
EXPOSE 8081