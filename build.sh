#!/bin/bash 
archs=(amd64 arm64) 
oss=(linux darwin windows)
for arch in ${archs[@]} 
do 
        for os in ${oss[@]}
        do
        env GOOS=${os} GOARCH=${arch} go build -o tagger_${os}_${arch} 
        chmod +x tagger_${os}_${arch} 
        done
done 