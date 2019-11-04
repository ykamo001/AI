.PHONY: setup

include build/makefiles/buildvars.mk

APPLICATION := ai
SERVICES := paint

build: setup
	go build \
	-ldflags "-X github.com/ykamo001/$(APPLICATION)/cmd.Buildstamp=$(BUILDSTAMP) -X github.com/ykamo001/$(APPLICATION)/cmd.Githash=$(GITHASH)"

test:
	go test ./... -v GOCACHE=off

integration_test:
	go test ./... -tags integration -v GOCACHE=off


setup: setup_vendor services

setup_vendor:
	@git config --global -l | grep -q 'url.git@github.com:.insteadof=https://github.com/' \
	|| (echo 'Update your git config for private repos with the following command.' \
		&& echo '(See https://albertech.blogpost.com/2016/11/fix-git-error-could-not-read-username.html)' \
		&& echo \
		&& echo 'git config --global --add url."git@github.com:".insteadOf "https://github.com"')
	go mod vendor

services:
	protoc --proto_path=$(GOPATH)/src:. --twirp_out=. --go_out=. ./rpc/paint/paint.proto;