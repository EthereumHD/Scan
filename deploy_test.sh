#!/bin/bash

HOST=gateway.poc.com
appname=scan
approot=/home/pocethereum/scan.service/

#./build.sh || exit 1
echo "start copy to remote"
scp scan.tar.gz pocethereum@$HOST:$approot
ssh pocethereum@$HOST "cd $approot && tar -xvzf scan.tar.gz && ./load.sh restart"

#
### done ############
echo "done,done,done"
