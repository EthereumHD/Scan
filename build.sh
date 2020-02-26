#!/usr/bin/env bash

#set -x
appname='scan'

###  build      ####
echo "start build ..."
#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -o ./bin/$appname ./src/main.go && (echo "BUILD SUCCESS"; exit 0;) || (echo "BUILD FAILED" && exit 1);
xgo -x -targets=linux/amd64  ./src/ && (mv src-linux-amd64 ./bin/$appname; echo "BUILD SUCCESS"; exit 0;) || (echo "BUILD FAILED" && exit 1) || exit 1;

dir=$(pwd)
echo "initalize ..."
echo "rm $dir/output"
rm -rf $dir/output/ 2>&1 > /dev/null
mkdir -p  output/bin
mkdir -p  output/conf
mkdir -p  output/log
mkdir -p  output/logs

### copy files  ####
echo "copy to destination dir"
mv ./bin/scan                  output/bin/$appname
cp ./conf/scan.conf.test      ./output/conf/


### shell script ####
cp ./load.sh ./output/

### tar ############
echo "tar ..."
cd output
tar -czf $appname.tar.gz ./bin ./conf  ./log ./load.sh
mv ./$appname.tar.gz $dir/
rm -r $dir/output/


