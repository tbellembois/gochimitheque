#!/usr/bin/env bash
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

cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/shema.sql /tmp/
cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/sample.sql /tmp/
cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/migration.sql /tmp/

cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/shema.sql .
cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/sample.sql .
cp /home/thbellem/workspace/workspace_rust/chimitheque_db/src/resources/migration.sql . 

cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit
docker compose up -d keycloak
export SQLITE_EXTENSION_DIR="/home/thbellem/workspace/workspace_rust/chimitheque_db/src/extensions"
cd /home/thbellem/workspace/workspace_rust/chimitheque_zmq_server || exit
RUST_LOG=debug cargo run -- --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/chimitheque.sqlite
# open SVG file a web browser!
#CARGO_PROFILE_RELEASE_DEBUG=true cargo flamegraph --dev -o /tmp/flamegraph.svg -- --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/chimitheque.sqlite
#RUST_LOG=error ./target/release/chimitheque_zmq_server --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/chimitheque.sqlite
