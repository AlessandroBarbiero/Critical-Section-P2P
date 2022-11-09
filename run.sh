#!/bin/bash

set -euo pipefail

filename="./clients.info"

while read line 
do 
   $TERM -e go run . "$line" &
done < "$filename"
