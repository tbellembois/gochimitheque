#!/usr/bin/env bash

debug=""
rustlog="INFO"

if [ "$CHIMITHEQUE_DEBUG" == "true" ]; then
    debug="-debug"
    rustlog="DEBUG"
fi

export SQLITE_EXTENSION_DIR="/var/www-data/extensions"
export DB_PATH="/data/chimitheque.sqlite"
export KEYCLOAK_BASE_URL=$CHIMITHEQUE_APPURL/keycloak
export KEYCLOAK_REALM="chimitheque"
export KEYCLOAK_CLIENT_ID="chimitheque"
export ADMINS=$CHIMITHEQUE_ADMINS
export RUST_LOG=$rustlog

command="/var/www-data/chimitheque_back"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -appurl=$CHIMITHEQUE_APPURL $debug"
echo "command:"
echo $command
$command
