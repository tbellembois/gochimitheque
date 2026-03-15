#!/usr/bin/env bash

trap 'kill 0' EXIT SIGINT SIGTERM

if [ "$1" == "-clean" ]; then
    cd /home/thbellem/workspace/workspace_rust/chimitheque_back
    cargo clean
    cd /home/thbellem/workspace/workspace_rust/chimitheque_db
    cargo clean
    cd /home/thbellem/workspace/workspace_rust/chimitheque_traits
    cargo clean
    cd /home/thbellem/workspace/workspace_rust/chimitheque_types
    cargo clean
    cd /home/thbellem/workspace/workspace_rust/chimitheque_utils
    cargo clean
fi

export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
export OTEL_TRACES_SAMPLER=parentbased_traceidratio
export OTEL_TRACES_SAMPLER_ARG=1.0
export RUST_LOG=chimitheque_back=info,chimitheque_db=debug,tower_http=warn
export DB_PATH="/home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/chimitheque.sqlite"
export KEYCLOAK_BASE_URL="https://192.168.1.18:8443/keycloak"
export KEYCLOAK_REALM="chimitheque"
export KEYCLOAK_CLIENT_ID="chimitheque"
export SQLITE_EXTENSION_DIR="/home/thbellem/workspace/workspace_rust/chimitheque_back/src/extensions"

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit

echo "#### starting keycloak ####"
docker compose up -d keycloak

echo "#### starting jaeger ####"
docker compose up -d jaeger

echo "#### starting frontend ####"
go run . &

cd /home/thbellem/workspace/workspace_rust/chimitheque_back || exit

echo "#### starting backend ####"
cargo run . &

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/dev || exit

echo "#### starting nginx ####"
./nginx-1.28.1-x86_64-linux &

wait
