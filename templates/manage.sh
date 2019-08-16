#!/bin/bash
set -xe

RS=${PWD##*/}
LINES="100"

while getopts ":s:p:v:l:" opt; do
  case $opt in
    s) SVC="$OPTARG"
    ;;
    p) PORT="$OPTARG"
    ;;
    v) SVCCMD="$OPTARG"
    ;;
    l) LINES="$OPTARG"
    ;;
    \?) echo "An invalid option has been entered: $OPTARG"
        echo
        exit 2
    ;;
    :)  echo "The additional argument for option $OPTARG was omitted."
        echo
        exit 3
    ;;
  esac
done

#Check Mandatory Options
if [ "s" == "x${SVC}" ]; then
  echo "-s service is required"
  echo
  exit 4
fi

#Check Mandatory Options
if [ "x" == "x${SVCCMD}" ]; then
  echo "-v service command is required"
  echo
  exit 4
fi

#Check Mandatory Options
if [ "x" == "x${PORT}" ]; then
  echo "-p port is required"
  echo
  exit 4
fi

ssh -p ${PORT} root@${RS} "systemctl ${SVCCMD} ${SVC}"

ssh -p ${PORT} root@${RS} "journalctl -u ${SVC} -n ${LINES}"


exit 0
