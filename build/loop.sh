#!/bin/sh

[ $# -ne 2 ] && echo "Usage: $0 <interval> <command>" && exit 1

trap 'exit 0;' SIGTERM

while true; do
  eval "$2"
  sleep "$1"
done
