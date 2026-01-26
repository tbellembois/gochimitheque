FROM golang:1.25-trixie AS builder
LABEL author="Thomas Bellembois"
ARG BuildID
ENV BuildID=${BuildID}

#
# Prepare.
#

# Install dependencies.
RUN apt -y update
RUN apt -y install openssl libssl-dev

# Create base directory.
RUN mkdir -p /home/thbellem/workspace \
    && ln -s /go /home/thbellem/workspace/workspace_go

# Install Jade.
RUN go install github.com/Joker/jade/cmd/jade@master

#
# Chimithèque Go sources.
#

# WASM: copy sources.
# sudo mount --bind ~/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-wasm ./bind-gochimitheque-wasm
WORKDIR /go/src/github.com/tbellembois/
COPY ./bind-gochimitheque-wasm ./gochimitheque-wasm

# BACKEND: copy sources.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
COPY . .

#
# Chimithèque Go build.
#

# WASM: build.
WORKDIR /go/src/github.com/tbellembois/gochimitheque-wasm
RUN GOOS=js GOARCH=wasm go get -v -d ./... \
    && GOOS=js GOARCH=wasm go build -o wasm .

# WASM: copy and compress binary.
RUN cp /go/src/github.com/tbellembois/gochimitheque-wasm/wasm /go/src/github.com/tbellembois/gochimitheque/wasm/ \
    && gzip -9 -v -c /go/src/github.com/tbellembois/gochimitheque/wasm/wasm > /go/src/github.com/tbellembois/gochimitheque/wasm/wasm.gz \
    && rm /go/src/github.com/tbellembois/gochimitheque/wasm/wasm

# FRONTEND: generate code.
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
RUN go generate

# FRONTEND: build.
RUN if [ -z $BuildID ]; then BuildID=$(date "+%Y%m%d"); fi; echo "BuildID=$BuildID"; go build -ldflags "-s -w -X main.BuildID=$BuildID"

#
# Chimithèque Rust sources.
#

# Install Rust.
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="$PATH:/root/.cargo/bin"

# Get sources.
WORKDIR /go/src/rust
RUN git clone https://github.com/tbellembois/chimitheque_back.git
RUN git clone https://github.com/tbellembois/chimitheque_db.git
RUN git clone https://github.com/tbellembois/chimitheque_types.git
RUN git clone https://github.com/tbellembois/chimitheque_traits.git
RUN git clone https://github.com/tbellembois/chimitheque_utils.git
RUN git clone https://github.com/tbellembois/chimitheque_pubchem.git

#
# Chimithèque Rust build.
#
WORKDIR /go/src/rust/chimitheque_back
RUN cargo build --release
RUN strip target/release/chimitheque_back

#
# Final image.
#

FROM golang:1.25-trixie

RUN apt -y update && apt -y upgrade
RUN update-ca-certificates -v

RUN rm -Rf /var/cache/apk

# Retrieve database shema sql files.
WORKDIR /tmp
RUN wget https://raw.githubusercontent.com/tbellembois/chimitheque_db/refs/heads/main/src/resources/shema.sql

# Ensure www-data user exists.
RUN addgroup --gid 82 --system chimitheque \
    && adduser --uid 82 --system --ingroup chimitheque chimitheque \
    && mkdir /data \
    && chown chimitheque /data \
    && chmod 700 /data \
    && mkdir /var/www-data \
    && mkdir /var/www-data/extensions \
    && chown chimitheque /var/www-data \
    && chown chimitheque /var/log \
    && chmod 755 /var/log

# Copy sql extensions.
COPY --from=builder /go/src/rust/chimitheque_back/src/extensions/ /var/www-data/extensions/

# Copy frontend binary.
COPY --from=builder /go/src/github.com/tbellembois/gochimitheque/gochimitheque /var/www-data/
RUN chown chimitheque /var/www-data/gochimitheque \
    && chmod +x /var/www-data/gochimitheque

# Copy backend binary.
COPY --from=builder /go/src/rust/chimitheque_back/target/release/chimitheque_back /var/www-data/
RUN chown chimitheque /var/www-data/chimitheque_back \
    && chmod +x /var/www-data/chimitheque_back

# Copying entrypoint.
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

# Container configuration.
USER chimitheque
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME [ "/data" ]
EXPOSE 8081
