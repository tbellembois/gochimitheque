#!/usr/bin/env bash

appurl=""
apppath=""
dockerport=""
enablepublicproductsendpoint=""
admins=""
logfile=""
debug=""
updateQRCode=""

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
if [ ! -z "$CHIMITHEQUE_LOGFILE" ]
then
      logfile="-logfile $CHIMITHEQUE_LOGFILE"
      echo $logfile
fi

if [ ! -z "$CHIMITHEQUE_UPDATEQRCODE" ]
then
      updateQRCode="-updateqrcode"
      echo $updateQRCode
fi

command="/var/www-data/chimitheque_utils_service"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -dbpath /data $listenport $appurl $apppath $dockerport $enablepublicproductsendpoint $admins $logfile $debug $updateQRCode"
echo "command:"
echo $command
$command