#!/usr/bin/env bash

appurl=""
apppath=""
dockerport=""
oidcdiscoverurl=""
oidcclientid=""
oidcclientsecret=""
admins=""
debug=""
rustlog="error"
updateqrcode=""

if [ ! -z "$CHIMITHEQUE_DOCKERPORT" ]; then
    dockerport="-dockerport $CHIMITHEQUE_DOCKERPORT"
fi
if [ ! -z "$CHIMITHEQUE_APPURL" ]; then
    appurl="-appurl $CHIMITHEQUE_APPURL"
fi
if [ ! -z "$CHIMITHEQUE_APPPATH" ]; then
    apppath="-apppath $CHIMITHEQUE_APPPATH"
fi
if [ ! -z "$CHIMITHEQUE_OIDCDISCOVERURL" ]; then
    oidcdiscoverurl="-oidcdiscoverurl $CHIMITHEQUE_OIDCDISCOVERURL"
fi
if [ ! -z "$CHIMITHEQUE_OIDCCLIENTID" ]; then
    oidcclientid="-oidcclientid $CHIMITHEQUE_OIDCCLIENTID"
fi
if [ ! -z "$CHIMITHEQUE_OIDCCLIENTSECRET" ]; then
    oidcclientsecret="-oidcclientsecret $CHIMITHEQUE_OIDCCLIENTSECRET"
fi
if [ ! -z "$CHIMITHEQUE_ADMINS" ]; then
    admins="-admins $CHIMITHEQUE_ADMINS"
fi
if [ "$CHIMITHEQUE_DEBUG" == "true" ]; then
    debug="-debug"
    rustlog="debug"
fi

export SQLITE_EXTENSION_DIR="/var/www-data/extensions"
export RUST_LOG=$rustlog

command="/var/www-data/chimitheque_zmq_server --db-path /data/chimitheque.sqlite"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -dbpath /data $appurl $apppath $dockerport $oidcdiscoverurl $oidcclientid $oidcclientsecret $admins $debug $updateqrcode"
echo "command:"
echo $command
$command
