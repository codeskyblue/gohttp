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
	sudo apt-get install -y nodejs npm

deps:
	npm install
	bower install

cross-build:
	GOOS=windows GOARCH=386 go build
	GOOS=linux GOARCH=386 go build -o fileserv-linux-386
	GOOS=linux GOARCH=amd64 go build -o fileserv-linux-amd64

webpack:
	(cd public; webpack)

# vim:ft=make
#
