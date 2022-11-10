#!/bin/bash

set -euo pipefail

filename="./clients.info"

echo "First node with token is $1"

while read line 
do 
   $TERM -e go run . "$line" $1 &
done < "$filename"
