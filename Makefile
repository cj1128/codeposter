flags := -X main.appVersion=`cat VERSION`
build_flags := $(flags) -s -w

install:
	godep go install -ldflags "$(flags)"
.PHONY: install

build:
	GOOS=darwin godep go build -ldflags "$(build_flags)" -o tmp/codeposter-oxs-$$(cat VERSION)
	GOOS=linux godep go build -ldflags "$(build_flags)" -o tmp/codeposter-linux-$$(cat VERSION)
	GOOS=windows godep go build -ldflags "$(build_flags)" -o tmp/codeposter-windows-$$(cat VERSION)
.PHONY: build
