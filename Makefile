.PHONY: test test-race lint

# client/ depends on ebiten/v2 which requires X11 dev headers (display server).
# Those packages are tested separately in an environment with a display.
GO_PACKAGES := $(shell go list ./... | grep -v /client)

test:
	go test $(GO_PACKAGES)

test-race:
	go test -race -count=1 $(GO_PACKAGES)

lint:
	golangci-lint run ./...
