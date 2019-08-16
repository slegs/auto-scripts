#!/bin/bash
set -xe

RS=${PWD##*/}

while getopts ":c:" OPT; do
  case ${OPT} in
    c) CMD=${OPTARG}
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
if [ "x" == "x${CMD}" ]; then
  echo "-c command is required"
  echo
  exit 4
fi

${CMD}

exit 0
