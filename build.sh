#!/bin/bash

# 设置目标平台
targets=(
    "linux/amd64"
    "linux/arm"
    "windows/amd64"
    "darwin/amd64"
)

# 创建 dist 目录
mkdir -p dist

# 循环编译每个目标平台
for target in "${targets[@]}"; do
    # 分割目标平台为 GOOS 和 GOARCH
    IFS="/" read -r GOOS GOARCH <<< "$target"

    # 设置输出文件名
    output="dist/jdcloudVm2Jumpserver-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output="${output}.exe"
    fi

    # 编译
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$output"

    if [ $? -ne 0 ]; then
        echo "Failed to build for $GOOS/$GOARCH"
        exit 1
    fi

    # 拷贝 config.yml 文件到对应目录
    cp config.yml "dist/config.yml"
    cp jd2jumpserver.sh "dist/jd2jumpserver.sh"

    # 单独拷贝linux版本
    cp jd2jumpserver.sh "dist/jd2jumpServer/jd2jumpserver.sh"
    cp config.yml "dist/jd2jumpServer/config.yml"
    cp dist/jdcloudVm2Jumpserver-linux-amd64 "dist/jd2jumpServer/jdcloudVm2Jumpserver"


    if [ $? -ne 0 ]; then
        echo "Failed to copy config.yml for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Build complete!"
