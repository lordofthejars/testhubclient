version ?= latest

build:
	go build -o testhubclient

cross:
	docker run -it --rm -v "$$PWD":/go/src/github.com/lordofthejars/testhubclient -w /go/src/github.com/lordofthejars/testhubclient -e "version=${version}" lordofthejars/goreleaser:1.0 crossbuild.sh