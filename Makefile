#
# Makefile
# hzsunshx, 2015-02-11 13:17
#

dev: webpack
	go build

prod: webpack
	(cd public; go-bindata -pkg public bundle.js css/)
	(cd templates; go-bindata -pkg templates ./...)
	go build -tags "bindata"

install-deps:
	sudo apt-get update -qq
	sudo apt-get install -qq nodejs npm

deps:
	npm install

cross-build:
	GOOS=windows GOARCH=386 go build
	GOOS=linux GOARCH=386 go build -o fileserv-linux-386
	GOOS=linux GOARCH=amd64 go build -o fileserv-linux-amd64

webpack:
	webpack

clean:
	rm public/bundle.js
# vim:ft=make
#
