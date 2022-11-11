@echo off
set filename=clients.info
set firstToken=%1

echo "First node with token is %firstToken%"

for /F "tokens=*" %%A in  (%filename%) do  (
   start cmd /k "go run main.go %%A %firstToken%"
)
@echo on
