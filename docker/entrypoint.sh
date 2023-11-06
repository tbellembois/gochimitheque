#!/usr/bin/env bash

appurl=""
apppath=""
dockerport=""
ldapserverurl=""
ldapserverusername=""
ldapserverpassword=""
ldapusersearchbasedn=""
ldapusersearchfilter=""
ldapgroupsearchbasedn=""
ldapgroupsearchfilter=""
autocreateuser=""
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

echo "parameters:"
if [ ! -z "$CHIMITHEQUE_LDAPSERVERURL" ]
then
      ldapserverurl="-ldapserverurl $CHIMITHEQUE_LDAPSERVERURL"
      echo $ldapserverurl
fi
if [ ! -z "$CHIMITHEQUE_LDAPSERVERUSERNAME" ]
then
      ldapserverusername="-ldapserverusername $CHIMITHEQUE_LDAPSERVERUSERNAME"
      echo $ldapserverusername
fi
if [ ! -z "$CHIMITHEQUE_LDAPSERVERPASSWORD" ]
then
      ldapserverpassword="-ldapserverpassword $CHIMITHEQUE_LDAPSERVERPASSWORD"
      echo $ldapserverpassword
fi
if [ ! -z "$CHIMITHEQUE_LDAPGROUPSEARCHBASEDN" ]
then
      ldapgroupsearchbasedn="-ldapgroupsearchbasedn $CHIMITHEQUE_LDAPGROUPSEARCHBASEDN"
      echo $ldapgroupsearchbasedn
fi
if [ ! -z "$CHIMITHEQUE_LDAPUSERSEARCHBASEDN" ]
then
      ldapusersearchbasedn="-ldapusersearchbasedn $CHIMITHEQUE_LDAPUSERSEARCHBASEDN"
      echo $ldapusersearchbasedn
fi
if [ ! -z "$CHIMITHEQUE_LDAPUSERSEARCHFILTER" ]
then
      ldapusersearchfilter="-ldapusersearchfilter $CHIMITHEQUE_LDAPUSERSEARCHFILTER"
      echo $ldapusersearchfilter
fi
if [ ! -z "$CHIMITHEQUE_LDAPGROUPSEARCHFILTER" ]
then
      ldapgroupsearchfilter="-ldapgroupsearchfilter $CHIMITHEQUE_LDAPGROUPSEARCHFILTER"
      echo $ldapgroupsearchfilter
fi
if [ ! -z "$CHIMITHEQUE_AUTOCREATEUSER" ]
then
      autocreateuser="-autocreateuser"
      echo $autocreateuser
fi
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
if [ ! -z "$CHIMITHEQUE_MAILSERVERADDRESS" ]
then
      mailserveraddress="-mailserveraddress $CHIMITHEQUE_MAILSERVERADDRESS"
      echo $mailserveraddress
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERPORT" ]
then
      mailserverport="-mailserverport $CHIMITHEQUE_MAILSERVERPORT"
      echo $mailserverport
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERSENDER" ]
then
      mailserversender="-mailserversender $CHIMITHEQUE_MAILSERVERSENDER"
      echo $mailserversender
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERUSETLS" ]
then
      mailserverusetls="-mailserverusetls"
      echo $mailserverusetls
fi
if [ ! -z "$CHIMITHEQUE_MAILSERVERTLSSKIPVERIFY" ]
then
      mailservertlsskipverify="-mailservertlsskipverify"
      echo $mailservertlsskipverify
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

if [ ! -z "$CHIMITHEQUE_RESETADMINPASSWORD" ]
then
      resetAdminPassword="-resetadminpassword"
      echo $resetAdminPassword
fi
if [ ! -z "$CHIMITHEQUE_UPDATEQRCODE" ]
then
      updateQRCode="-updateqrcode"
      echo $updateQRCode
fi
if [ ! -z "$CHIMITHEQUE_MAILTEST" ]
then
      mailTest="-mailtest $CHIMITHEQUE_MAILTEST"
      echo $mailTest
fi
if [ ! -z "$CHIMITHEQUE_IMPORTFROM" ]
then
      importfrom="-importfrom $CHIMITHEQUE_IMPORTFROM"
      echo $importfrom
fi

command="/var/www-data/chimitheque_utils_service"
echo "command:"
echo $command
$command &

command="/var/www-data/gochimitheque -dbpath /data $listenport $appurl $apppath $dockerport $ldapserverurl $ldapserverusername $ldapserverpassword $ldapgroupsearchbasedn $ldapgroupsearchfilter $ldapusersearchbasedn $ldapusersearchfilter $autocreateuser $mailserveraddress $mailserverport $mailserversender $mailserverusetls $mailservertlsskipverify $enablepublicproductsendpoint $admins $logfile $debug $resetAdminPassword $updateQRCode $mailTest $importfrom"
echo "command:"
echo $command
$command