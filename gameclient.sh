#!/bin/bash
cd "$(dirname "$0")"
go build -o bin/client src/client/main.go
chmod 777 bin/client
touch /var/log/battle.log
chmod 666 /var/log/battle.log
bin/client $1 $2 $3 2>/var/log/battle.log
