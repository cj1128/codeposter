SHELL := /bin/bash

flags := -X main.appVersion=`git tag --points-at HEAD`

install:
	go install -ldflags "$(flags)"
.PHONY: install

build:
	@if [[ $$(git tag --points-at HEAD) == "" ]]; then \
		echo 'need a tag to build!' ; \
		exit 1 ; \
	fi

	# windows
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -tags static -ldflags "-s -w" -o "tmp/codeposter_windows_amd64.exe"

	# mac
	go build -v -tags static -ldflags "-s -w" -o "tmp/codeposter_darwin_amd64"
.PHONY: build
