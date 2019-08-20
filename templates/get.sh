#!/bin/bash

set -xe

RS=${PWD##*/}
ABSOLUTE=${PWD}
NOW="$(date +%Y%m%d%H%M%S)"
FILTER=""

while getopts ":d:b:f:p:t:e:i:" opt; do
  case $opt in
    d) DIR1="${OPTARG}"
    ;;
    b) BACKUPDIR="${OPTARG}"
    ;;
    f) FILESDIR="${OPTARG}"
    ;;
    p) PORT="${OPTARG}"
    ;;
    t) NOW="${OPTARG}"
    ;;
    i) FILTER="${OPTARG}"
    ;;
    e) EXCEPT="${OPTARG}"
    ;;
    \?) echo "An invalid option has been entered: ${OPTARG}"
        echo
        exit 2
    ;;
    :)  echo "The additional argument for option ${OPTARG} was omitted."
        echo
        exit 3
    ;;
  esac
done

#Check Mandatory Options
if [ "x" == "x${DIR1}}" ]; then
  echo "-d target directory is required"
  echo
  exit 4
fi

#Check Mandatory Options
if [ "x" == "x${PORT}}" ]; then
  echo "-p port is required"
  echo
  exit 4
fi

if [ "x" == "x${BACKUPDIR}}" ]; then
  echo "-b backup directory is required"
  echo
  exit 4
fi

if [ "x" == "x${FILESDIR}}" ]; then
  echo "-f directory for storing files is required"
  echo
  exit 4
fi

#Deal with excludes
EXCLUDESTRING=""
if [ "x" != "x${EXCEPT}}" ]; then
  for I in $(echo ${EXCEPT} | sed "s/,/ /g")
  do
      # call your procedure/other scripts here below
      EXCLUDESTRING="${EXCLUDESTRING} --exclude ${I} "
  done
fi

mkdir -p ${FILESDIR}${DIR1}

#./delete-backups.sh "${BACKUPDIR}"

sleep 1

rsync -rlptDiv -e "ssh -p ${PORT}" --progress --backup --backup-dir=${ABSOLUTE}${BACKUPDIR}/${NOW}${DIR1} --suffix=".deleted" --delete ${EXCLUDESTRING} root@${RS}:${DIR1}/${FILTER} ${FILESDIR}${DIR1}/${FILTER}

exit 0
