#!/usr/bin/env bash

appurl=""
apppath=""
dockerport=""
oidcdiscoverurl=""
oidcclientid=""
oidcclientsecret=""
admins=""
debug=""
updateqrcode=""

echo "parameters:"

if [ ! -z "$CHIMITHEQUE_DOCKERPORT" ]
then
      dockerport="-dockerport $CHIMITHEQUE_DOCKERPORT"
      echo $dockerport
fi
if [ ! -z "$CHIMITHEQUE_APPURL" ]
then
      appurl="-appurl $CHIMITHEQUE_APPURL"
      echo $appurl
fi
if [ ! -z "$CHIMITHEQUE_APPPATH" ]
then
      apppath="-apppath $CHIMITHEQUE_APPPATH"
      echo $apppath
fi
if [ ! -z "$CHIMITHEQUE_OIDCDISCOVERURL" ]
then
      oidcdiscoverurl="-oidcdiscoverurl $CHIMITHEQUE_OIDCDISCOVERURL"
      echo $oidcdiscoverurl
fi
if [ ! -z "$CHIMITHEQUE_OIDCCLIENTID" ]
then
      oidcclientid="-oidcclientid $CHIMITHEQUE_OIDCCLIENTID"
      echo $oidcclientid
fi
if [ ! -z "$CHIMITHEQUE_OIDCCLIENTSECRET" ]
then
      oidcclientsecret="-oidcclientsecret $CHIMITHEQUE_OIDCCLIENTSECRET"
      echo $oidcclientsecret
fi

if [ ! -z "$CHIMITHEQUE_ADMINS" ]
then
      admins="-admins $CHIMITHEQUE_ADMINS"
      echo $admins
fi
if [ ! -z "$CHIMITHEQUE_DEBUG" ]
then
      debug="-debug"
      echo $debug
fi

if [ ! -z "$CHIMITHEQUE_UPDATEQRCODE" ]
then
      updateqrcode="-updateqrcode"
      echo $updateqrcode
fi

export SQLITE_EXTENSION_DIR="/var/www-data/extensions"

command="/var/www-data/chimitheque_zmq_server --db-path /data/storage.db"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -dbpath /data $appurl $apppath $dockerport $oidcdiscoverurl $oidcclientid $oidcclientsecret $admins $debug $updateqrcode"
echo "command:"
echo $command
$command
