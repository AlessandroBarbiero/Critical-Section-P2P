#!/bin/bash

set -euo pipefail

filename = "./clients.info"

while read line 
do 
    go run . "$line" &
done < filename
