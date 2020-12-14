#!/usr/bin/env bash
proxyurl=""
proxypath=""
mailserveraddress=""
mailserverport=""
mailserversender=""
mailserverusetls=""
mailservertlsskipverify=""
enablepublicproductsendpoint=""
admins=""
logfile=""
debug=""

resetAdminPassword=""
updateQRCode=""
mailTest=""
importv1from=""
importfrom=""

if [ ! -z "$CHIMITHEQUE_PROXYURL" ]
then
      proxyurl="-proxyurl $CHIMITHEQUE_PROXYURL"
fi
if [ ! -z "$CHIMITHEQUE_PROXYPATH" ]
then
      proxypath="-proxypath $CHIMITHEQUE_PROXYPATH"
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERADDRESS" ]
then
      mailserveraddress="-mailserveraddress $CHIMITHEQUE_MAILSERVERADDRESS"
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERPORT" ]
then
      mailserverport="-mailserverport $CHIMITHEQUE_MAILSERVERPORT"
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERSENDER" ]
then
      mailserversender="-mailserversender $CHIMITHEQUE_MAILSERVERSENDER"
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERUSETLS" ]
then
      mailserverusetls="-mailserverusetls"
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERTLSSKIPVERIFY" ]
then
      mailservertlsskipverify="-mailservertlsskipverify"
fi
if [ ! -z "$CHIMITHEQUE_ENABLEPUBLICPRODUCTSENDPOINT" ]
then
      enablepublicproductsendpoint="-enablepublicproductsendpoint"
fi
if [ ! -z "$CHIMITHEQUE_ADMINS" ]
then
      admins="-admins $CHIMITHEQUE_ADMINS"
fi
if [ ! -z "$CHIMITHEQUE_DEBUG" ]
then
      debug="-debug"
fi
if [ ! -z "$CHIMITHEQUE_LOGFILE" ]
then
      logfile="-logfile $CHIMITHEQUE_LOGFILE"
fi

if [ ! -z "$CHIMITHEQUE_RESETADMINPASSWORD" ]
then
      resetAdminPassword="-resetadminpassword"
fi
if [ ! -z "$CHIMITHEQUE_UPDATEQRCODE" ]
then
      updateQRCode="-updateqrcode"
fi
if [ ! -z "$CHIMITHEQUE_MAILTEST" ]
then
      mailTest="-mailtest"
fi
if [ ! -z "$CHIMITHEQUE_IMPORTV1FROM" ]
then
      importv1from="-importv1from"
fi
if [ ! -z "$CHIMITHEQUE_IMPORTFROM" ]
then
      importfrom="-importfrom $CHIMITHEQUE_IMPORTFROM"
fi

/var/www-data/gochimitheque -dbpath /data \
$listenport \
$proxyurl \
$proxypath \
$mailserveraddress \
$mailserverport \
$mailserversender \
$mailserverusetls \
$mailservertlsskipverify \
$enablepublicproductsendpoint \
$admins \
$logfile \
$debug \
$resetAdminPassword \
$updateQRCode \
$mailTest \
$importv1from \
$importfrom