export GO111MODULE := on
GITHASH := $(shell git rev-parse HEAD)
BUILDSTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
GITHUB_URL := $(shell git remote -v | awk '/^origin[\t ]git@github.com:(.*)\.git \(fetch\)/ { split($$0, arr, "(origin\tgit@github.com:)|(.git)"); print("https://github.com/" arr[2]); }')
COMMIT_URL := $(GITHUB_URL)/commit/$(GITHASH)