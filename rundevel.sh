#!/usr/bin/env bash
cd /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque || exit
docker compose up -d keycloak
cd /home/thbellem/workspace/workspace_rust/chimitheque_utils_service  || exit
RUST_LOG=debug cargo run -- --db-path /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque/storage.db