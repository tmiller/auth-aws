appname := auth-aws
sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(appname)

tar = cd build && tar -czvf $(appname)_$(1)_$(2).tgz $(appname) && rm $(appname)

.PHONY: all clean linux darwin

all: linux darwin

clean:
	rm -rf build

linux: build/auth-aws_linux_386.tgz build/auth-aws_linux_amd64.tgz

build/auth-aws_linux_386.tgz: $(sources)
	$(call build,linux,386)
	$(call tar,linux,386)

build/auth-aws_linux_amd64.tgz: $(sources)
	$(call build,linux,amd64)
	$(call tar,linux,amd64)

darwin:  build/auth-aws_darwin_386.tgz build/auth-aws_darwin_amd64.tgz

build/auth-aws_darwin_386.tgz: $(sources)
	$(call build,darwin,386)
	$(call tar,darwin,386)

build/auth-aws_darwin_amd64.tgz: $(sources)
	$(call build,darwin,amd64)
	$(call tar,darwin,amd64)
