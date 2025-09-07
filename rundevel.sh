#!/usr/bin/env bash
cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit
docker compose up -d keycloak
export SQLITE_EXTENSION_DIR="/home/thbellem/workspace/workspace_rust/chimitheque_db/src/extensions"
cd /home/thbellem/workspace/workspace_rust/chimitheque_zmq_server || exit
RUST_LOG=debug cargo run -- --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/storage.db
# open SVG file a web browser!
#CARGO_PROFILE_RELEASE_DEBUG=true cargo flamegraph --dev -o /tmp/flamegraph.svg -- --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/storage.db
#RUST_LOG=error ./target/release/chimitheque_zmq_server --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/storage.db
