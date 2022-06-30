#!/usr/bin/env bash

# 针对非go module的项目，需要指定目录
#cpath=`pwd`
#PROJECT_PATH=${cpath%src*}
#echo $PROJECT_PATH
#export GOPATH=$GOPATH:${PROJECT_PATH}


TARGET_PATH=./bin
SOURCE_PATH=./brun
TARGET_FILE_NAME=reskd
SOURCE_FILE_NAME=main

rm -fr ${TARGET_FILE_NAME}*

build(){
    echo $GOOS $GOARCH
    tname=${TARGET_FILE_NAME}_${GOOS}_${GOARCH}${EXT}
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ${TARGET_PATH}/${tname} -v ${SOURCE_PATH}/${SOURCE_FILE_NAME}.go

    chmod +x ${TARGET_PATH}/${tname}

    mv ${TARGET_PATH}/${tname} ${TARGET_PATH}/${TARGET_FILE_NAME}${EXT}
    if [ ${GOOS} == "windows" ];then
        zip ${TARGET_PATH}/${tname}.zip ${TARGET_PATH}/${TARGET_FILE_NAME}${EXT} config.ini ../public/
    else
        tar --exclude=*.gz  --exclude=*.zip  --exclude=*.git -czvf ${TARGET_PATH}/${tname}.tar.gz ${TARGET_PATH}/${TARGET_FILE_NAME}${EXT} config.ini *.sh ../public/ -C ./ .
    fi
    mv ${TARGET_PATH}/${TARGET_FILE_NAME}${EXT} ${TARGET_PATH}/${tname}

}
CGO_ENABLED=0

#mac os 64
GOOS=darwin
GOARCH=amd64
build

##linux 64
#GOOS=linux
#GOARCH=amd64
#build
#
##windows
##64
#EXT=.exe
#GOOS=windows
#GOARCH=amd64
#build


