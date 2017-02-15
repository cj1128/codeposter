.PHONYY: install

install:
	godep go install -ldflags "-X main.appVersion=`cat VERSION`"
