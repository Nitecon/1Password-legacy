#!/usr/bin/env bash
APP_NAME="1pass"
SERVER_BIN="svx.bin"
EVERS="0.36.7"
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
appDir="$( cd -P "$( dirname "$SOURCE" )/" && pwd )"

buildLinux64(){
    echo "-> Building linux/amd64 server thread"
    mkdir -p .build/$1/$2
    GOOS=linux GOARCH=amd64 go build -o .build/linux/amd64/svx.bin .

    echo "-> Building linux/x86_64 application binary"
    if [ -f "$appDir/app/$SERVER_BIN" ]; then
        rm -f "$appDir/app/$SERVER_BIN"
    fi
    mv $appDir/.build/linux/amd64/svx.bin $appDir/app/
    cd dist
    $appDir/node_modules/electron-packager/cli.js $appDir/app $APP_NAME --platform=linux --arch=x64 --version=$EVERS
    rm -rf $appDir/.build/
    echo "Linux 64bit Binary built: $appDir/dist/1pass-linux-x64/1pass"
}

cd $appDir
echo "-->> Cleaning out build dir"
rm -rf .build/
echo "-->> Cleaning out dist dir"
rm -rf dist/
mkdir -p dist
echo "-->> Installing dependencies"
npm update
echo "-->> Upgrading all modules"
npm update electron-prebuilt

echo "--> Building server thread binaries"
cd $appDir
echo "--> Building Linux 64"
buildLinux64
