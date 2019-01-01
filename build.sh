#!/bin/sh

cd `dirname $0`
rm -rf bin
mkdir bin

for ARCH in 386 amd64 ; do
	OS=windows
	GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-s -w"
	zip -r "bin/${OS}_${ARCH}.zip" oogrep.exe
	rm oogrep.exe

	for OS in darwin linux; do
		GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-s -w"
		tar -zcvf "bin/${OS}_${ARCH}.tgz" oogrep
	done
done

