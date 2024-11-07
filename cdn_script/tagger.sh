#!/bin/sh

bsae="http://get.mereith.com/tagger/"

# 从互联网安装 tagger 脚本
# 检查系统是 amd64 还是 arm64
if [ $(uname -m) = "x86_64" ]; then
    arch="amd64"
else
    arch="arm64"
fi

# 检查系统 darwin windows linux
if [ $(uname) = "Darwin" ]; then
    os="darwin"
elif [ $(uname) = "Linux" ]; then
    os="linux"
else
    os="windows"
fi

# windows 不复制
if [ $os = "windows" ]; then
    echo "windows 不支持安装"
    exit 1
fi

tagger="tagger_${os}_${arch}"


echo 系统为 ${os} 架构为 ${arch}\n\n

# 使用 curl 下载文件
curl -L $bsae$tagger -o tagger



# 检查是否有 sudo
if [ -x "$(command -v sudo)" ]; then
    # 删除老版本的 tagger
    sudo rm -f /usr/local/bin/tagger
    sudo rm -f /usr/local/sbin/tagger
    sudo rm -rf ~/.tagger
    sudo chmod +x ./tagger
    sudo cp ./tagger /usr/bin/tagger
else
    rm -f /usr/local/bin/tagger
    rm -f /usr/local/sbin/tagger
    rm -rf ~/.tagger
    chmod +x ./tagger
    cp ./tagger  /usr/bin/tagger
fi
echo 安装完成！