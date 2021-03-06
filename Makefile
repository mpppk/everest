SHELL = /bin/bash

.PHONY: lint
lint: generate
	go vet $(go list ./... | grep -v /scripts/)

.PHONY: test
test: generate
	go test ./...

.PHONY: integration-test
integration-test:
	go test -tags=integration ./...

.PHONY: coverage
coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov:  coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: clean
clean:
	rm -rf self

.PHONY: build-front
build-front:
	cd defaultembedded; yarn install
	cd defaultembedded; yarn build

.PHONY: generate
generate: clean build-front
	go generate

.PHONY: build
build: generate
	go build

.PHONY: build-macos-app
build-macos-app: generate
	go run scripts/build-mac/build-mac

.PHONY: cross-build-snapshot
cross-build:
	goreleaser --rm-dist --snapshot

.PHONY: install
install: generate
	go install

.PHONY: circleci
circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN

.PHONY: macos-icon
macos-icon:
	iconutil -c icns everest.iconset
	mv -f everest.icns defaultembedded/src/
