#!/bin/bash

# export GOPATH="${PWD}"

# GOARCH="arm"

# GOOS=android
# TOOLCHAIN_DIR=./pkg/gomobile/ndk-toolchains
# # for GOARM in 7 6 5; do
# 	CC="${TOOLCHAIN_DIR}/arm/bin/clang"
# 	CXX="${TOOLCHAIN_DIR}/arm/bin/clang"
# 	CC_FOR_TARGET="${TOOLCHAIN_DIR}/arm/bin/clang"
# 	GOARM="7"
# 	echo $CC
# 	echo $CXX
# 	echo $CC_FOR_TARGET

# 	BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}${GOARM}"
# 	CMD="CC_FOR_TARGE=${CC_FOR_TARGET} GOARM=${GOARM} GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} main"
# 	echo "${CMD}"
#     eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}" 
# # done

export GOPATH=`pwd`
export CGO_ENABLED=1
export GOOS=android 
export LDFLAGS="-s -w"

export NDK_TOOLCHAIN=/home/yourok/Space/Projects/Android/TorrServe/TorrServer/pkg/gomobile/ndk-toolchains/arm
export CC=$NDK_TOOLCHAIN/bin/arm-linux-androideabi-clang
export CXX=$NDK_TOOLCHAIN/bin/arm-linux-androideabi-clang++
export GOARCH=arm 
export GOARM=7
BIN_FILENAME="dist/TorrServer-${GOOS}-${GOARCH}${GOARM}"
echo "Android ${BIN_FILENAME}"
go build -ldflags="${LDFLAGS}" -o ${BIN_FILENAME} main

export NDK_TOOLCHAIN=/home/yourok/Space/Projects/Android/TorrServe/TorrServer/pkg/gomobile/ndk-toolchains/arm64
export CC=$NDK_TOOLCHAIN/bin/aarch64-linux-android-clang
export CXX=$NDK_TOOLCHAIN/bin/aarch64-linux-android-clang++
export GOARCH=arm64
export GOARM=""
BIN_FILENAME="dist/TorrServer-${GOOS}-${GOARCH}${GOARM}"
echo "Android ${BIN_FILENAME}"
go build -ldflags="${LDFLAGS}" -o ${BIN_FILENAME} main

export NDK_TOOLCHAIN=/home/yourok/Space/Projects/Android/TorrServe/TorrServer/pkg/gomobile/ndk-toolchains/x86
export CC=$NDK_TOOLCHAIN/bin/i686-linux-android-clang
export CXX=$NDK_TOOLCHAIN/bin/i686-linux-android-clang++
export GOARCH=386
export GOARM=""
BIN_FILENAME="dist/TorrServer-${GOOS}-${GOARCH}${GOARM}"
echo "Android ${BIN_FILENAME}"
go build -ldflags="${LDFLAGS}" -o ${BIN_FILENAME} main

export NDK_TOOLCHAIN=/home/yourok/Space/Projects/Android/TorrServe/TorrServer/pkg/gomobile/ndk-toolchains/x86_64
export CC=$NDK_TOOLCHAIN/bin/x86_64-linux-android-clang
export CXX=$NDK_TOOLCHAIN/bin/x86_64-linux-android-clang++
export GOARCH=amd64
export GOARM=""
BIN_FILENAME="dist/TorrServer-${GOOS}-${GOARCH}${GOARM}"
echo "Android ${BIN_FILENAME}"
go build -ldflags="${LDFLAGS}" -o ${BIN_FILENAME} main