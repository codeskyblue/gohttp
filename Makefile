#
# Makefile
# hzsunshx, 2015-02-11 13:17
#

everything: npm webpack binary

binary:
	#go-bindata data templates public/...
	go build

npm:
	npm install

cross-build:
	GOOS=windows GOARCH=386 go build
	GOOS=linux GOARCH=386 go build -o fileserv-linux-386
	GOOS=linux GOARCH=amd64 go build -o fileserv-linux-amd64

webpack:
	(cd public; webpack)

# vim:ft=make
#
