#!/usr/bin/env bash

appurl=""
apppath=""
dockerport=""
oidcissuer=""
oidcclientid=""
oidcclientsecret=""
enablepublicproductsendpoint=""
admins=""
debug=""
updateqrcode=""
autoimporturl=""

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
if [ ! -z "$CHIMITHEQUE_OIDCISSUER" ]
then
      oidcissuer="-oidcissuer $CHIMITHEQUE_OIDCISSUER"
      echo $oidcissuer
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

if [ ! -z "$CHIMITHEQUE_ENABLEPUBLICPRODUCTSENDPOINT" ]
then
      enablepublicproductsendpoint="-enablepublicproductsendpoint"
      echo $enablepublicproductsendpoint
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
if [ ! -z "$CHIMITHEQUE_AUTOIMPORT_URL" ]
then
      autoimporturl="-autoimporturl $CHIMITHEQUE_AUTOIMPORT_URL"
      echo $autoimporturl
fi

if [ ! -z "$CHIMITHEQUE_UPDATEQRCODE" ]
then
      updateqrcode="-updateqrcode"
      echo $updateqrcode
fi

command="/var/www-data/chimitheque_utils_service"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -dbpath /data $listenport $appurl $apppath $dockerport $oidcissuer $oidcclientid $oidcclientsecret $enablepublicproductsendpoint $admins $debug $autoimporturl $updateqrcode"
echo "command:"
echo $command
$command