# Stage 1: Build Go components
FROM golang:1.25-trixie AS go_builder

LABEL author="Thomas Bellembois"
ARG BuildID

#
# Prepare.
#
# Install dependencies.
RUN apt -y update && \
    apt -y install openssl libssl-dev

# Create base directory.
RUN mkdir -p /home/thbellem/workspace \
    && ln -s /go /home/thbellem/workspace/workspace_go

# Install Jade.
RUN go install github.com/Joker/jade/cmd/jade@master

#
# Chimithèque Go sources.
#

# Copy sources and build
# sudo mount --bind ~/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-wasm ./bind-gochimitheque-wasm
WORKDIR /go/src/github.com/tbellembois/
COPY ./bind-gochimitheque-wasm ./gochimitheque-wasm
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
COPY . .

WORKDIR /go/src/github.com/tbellembois/gochimitheque-wasm
RUN GOOS=js GOARCH=wasm go get -v ./... && \
    GOOS=js GOARCH=wasm go build -o wasm

# Compress binary
RUN cp /go/src/github.com/tbellembois/gochimitheque-wasm/wasm /go/src/github.com/tbellembois/gochimitheque/wasm/ && \
    gzip -9 -v -c /go/src/github.com/tbellembois/gochimitheque-wasm/wasm > /go/src/github.com/tbellembois/gochimitheque/wasm/wasm.gz && \
    rm /go/src/github.com/tbellembois/gochimitheque/wasm/wasm

# Build backend
WORKDIR /go/src/github.com/tbellembois/gochimitheque/
RUN go generate
RUN if [ -z $BuildID ]; then BuildID=$(date "+%Y%m%d"); fi; echo "BuildID=$BuildID"; \
    go build -ldflags "-s -w -X main.BuildID=$BuildID"

# Stage 2: Build Rust components
FROM rust:1.93.1 AS rust_builder

WORKDIR /go/src/rust
RUN git clone https://github.com/tbellembois/chimitheque_back.git && \
    git clone https://github.com/tbellembois/chimitheque_db.git && \
    git clone https://github.com/tbellembois/chimitheque_types.git && \
    git clone https://github.com/tbellembois/chimitheque_traits.git && \
    git clone https://github.com/tbellembois/chimitheque_utils.git && \
    git clone https://github.com/tbellembois/chimitheque_pubchem.git

WORKDIR /go/src/rust/chimitheque_back
RUN cargo build --release --locked && \
    strip target/release/chimitheque_back

# Stage 3: Final image
FROM golang:1.25-trixie AS final_image

RUN apt -y update && \
    apt -y upgrade && \
    update-ca-certificates -v

RUN rm -Rf /var/cache/apk

# WORKDIR /tmp
# RUN wget https://raw.githubusercontent.com/tbellembois/chimitheque_db/refs/heads/main/src/resources/shema.sql

# Ensure www-data user exists.
RUN addgroup --gid 82 --system chimitheque && \
    adduser --uid 82 --system --ingroup chimitheque chimitheque && \
    mkdir /data && \
    chown chimitheque /data && \
    chmod 700 /data && \
    mkdir -p /var/www-data/extensions && \
    chown chimitheque /var/www-data/extensions && \
    chown chimitheque /var/log && \
    chmod 755 /var/log

# Copy SQL extensions.
COPY --from=rust_builder /go/src/rust/chimitheque_back/src/extensions/ /var/www-data/extensions/

# Copy frontend binary.
COPY --from=go_builder /go/src/github.com/tbellembois/gochimitheque/wasm/wasm.gz /var/www-data/
COPY --from=go_builder /go/src/github.com/tbellembois/gochimitheque/gochimitheque /var/www-data/
# RUN chown chimitheque /var/www-data/{gochimitheque,wasm.gz} && \
#     chmod +x /var/www-data/{gochimitheque,wasm.gz}
RUN chown chimitheque /var/www-data/gochimitheque && \
    chown chimitheque /var/www-data/wasm.gz && \
    chmod +x /var/www-data/gochimitheque

# Copy backend binary.
COPY --from=rust_builder /go/src/rust/chimitheque_back/target/release/chimitheque_back /var/www-data/
RUN chown chimitheque /var/www-data/chimitheque_back && \
    chmod +x /var/www-data/chimitheque_back

# Copying entrypoint.
COPY docker/entrypoint.sh /
RUN chmod +x /entrypoint.sh

USER chimitheque
WORKDIR /var/www-data
ENTRYPOINT [ "/entrypoint.sh" ]
VOLUME [ "/data" ]
EXPOSE 8081
