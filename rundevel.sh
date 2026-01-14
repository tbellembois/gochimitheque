#!/usr/bin/env bash

trap 'kill 0' EXIT SIGINT SIGTERM

if [ "$1" == "-clean" ]; then
    cd /home/thbellem/workspace/workspace_rust/chimitheque_zmq_server
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

echo "#### starting keycloak ####"

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit
docker compose up -d keycloak

echo "#### starting backend ####"

cd /home/thbellem/workspace/workspace_rust/chimitheque_back || exit
export SQLITE_EXTENSION_DIR="/home/thbellem/workspace/workspace_rust/chimitheque_back/src/extensions"
RUST_LOG="debug" DB_PATH="/home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/chimitheque.sqlite" KEYCLOAK_BASE_URL="https://192.168.1.18:8443/keycloak" KEYCLOAK_REALM="chimitheque" CLIENT_ID="chimitheque" cargo run . &

echo "#### starting frontend ####"

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit
go run . &

echo "#### starting nginx ####"

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/dev || exit
./nginx-1.28.1-x86_64-linux &

wait
