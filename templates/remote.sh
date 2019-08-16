#!/bin/bash
set -xe

RS=${PWD##*/}

while getopts ":c:p:" opt; do
  case $opt in
    c) CMD="$OPTARG"
    ;;
    p) PORT="$OPTARG"
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
if [ "x" == "x${CMD}" ]; then
  echo "-c command is required"
  echo
  exit 4
fi

#Check Mandatory Options
if [ "x" == "x${PORT}" ]; then
  echo "-p port is required"
  echo
  exit 4
fi

ssh -p ${PORT} root@${RS} "${CMD}"


exit 0
