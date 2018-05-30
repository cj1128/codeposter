flags := -X main.appVersion=`git tag --points-at HEAD`
build_flags := $(flags) -s -w

install:
	go install -ldflags "$(flags)"
.PHONY: install

build:
	gox -arch="386 amd64" -os="darwin linux windows" -ldflags="$(build_flags)" -output="tmp/{{.Dir}}_{{.OS}}_{{.Arch}}"
.PHONY: build
