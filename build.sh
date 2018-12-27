#!/bin/sh

for ARCH in 386 amd64 ; do
	OS=windows
	mkdir -p "bin/$OS/$ARCH"
	GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-s -w"
	mv oogrep.exe "bin/$OS/$ARCH"
	zip -r "bin/${OS}_${ARCH}.zip" "bin/$OS/$ARCH"

	for OS in darwin linux; do
		mkdir -p "bin/$OS/$ARCH"
		GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-s -w"
		mv oogrep "bin/$OS/$ARCH"
		tar -zcvf "bin/${OS}_${ARCH}.tgz" "bin/$OS/$ARCH"
	done
done

