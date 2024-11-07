#!/bin/sh

base="http://get.mereith.com/tagger/"

# 从互联网安装 tagger 脚本
# 检查系统是 amd64 还是 arm64
if [ $(uname -m) = "x86_64" ]; then
    arch="amd64"
elif [ $(uname -m) = "aarch64" ] || [ $(uname -m) = "arm64" ]; then
    arch="arm64"
else
    echo "不支持的架构: $(uname -m)"
    exit 1
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
echo "正在下载 $tagger ..."
if ! curl -L $base$tagger -o tagger; then
    echo "下载失败"
    exit 1
fi

# 检查是否有 sudo
if [ -x "$(command -v sudo)" ]; then
    # 删除老版本的 tagger
    sudo rm -f /usr/local/bin/tagger
    sudo rm -f /usr/local/sbin/tagger
    sudo rm -rf ~/.tagger
    
    if [ ! -d /usr/local/bin ]; then
        sudo mkdir -p /usr/local/bin
    fi
    
    sudo chmod +x ./tagger
    sudo cp ./tagger /usr/local/bin/tagger
else
    rm -f /usr/local/bin/tagger
    rm -f /usr/local/sbin/tagger
    rm -rf ~/.tagger
    
    if [ ! -d /usr/local/bin ]; then
        mkdir -p /usr/local/bin
    fi
    
    chmod +x ./tagger
    cp ./tagger /usr/local/bin/tagger
fi

# 清理临时文件
rm -f ./tagger

echo "安装完成！"