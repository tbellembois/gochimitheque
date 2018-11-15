#!/bin/bash
rm dist/*

echo "generating rice"
go generate

echo "building windows10 .exe"
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o dist/gochimitheque-win10.exe
echo "building windows7 .exe"
CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -o dist/gochimitheque-win7.exe
echo "building linux binary"
go build -o dist/gochimitheque-linux

rsync -av sample_* /tmp/
rsync -av dist/* /tmp

echo "building demo windows10 zip"
zip dist/gochimitheque-win10-demo.zip /tmp/sample_* /tmp/gochimitheque-win10.exe
echo "building demo windows7 zip"
zip dist/gochimitheque-win7-demo.zip /tmp/sample_* /tmp/gochimitheque-win7.exe
echo "building demo linux zip"
zip dist/gochimitheque-linux-demo.zip /tmp/sample_* /tmp/gochimitheque-linux