SHELL := /bin/bash

flags := -X main.appVersion=`git tag --points-at HEAD`
build_flags := $(flags) -s -w

install:
	go install -ldflags "$(flags)"
.PHONY: install

build:
	@if [[ $$(git tag --points-at HEAD) == "" ]]; then \
		echo 'need a tag to build!' ; \
		exit 1 ; \
	fi

	go-bindata -prefix static static

	gox -arch="386 amd64" -os="darwin linux windows" -ldflags="$(build_flags)" -output="tmp/{{.Dir}}_{{.OS}}_{{.Arch}}"
.PHONY: build
