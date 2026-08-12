.PHONY: dev run build cli test gui-build

ARGS ?=
ROOT_DIR := $(CURDIR)
GUI_DIR := $(ROOT_DIR)/packages/gui
GO_CACHE ?= /tmp/hivelab-go-build

dev:
	cd $(GUI_DIR) && GOCACHE=$(GO_CACHE) wails dev $(ARGS)

run: dev

build:
	cd $(GUI_DIR) && GOCACHE=$(GO_CACHE) wails build $(ARGS)

cli:
	GOCACHE=$(GO_CACHE) go run ./packages/cli $(ARGS)

test:
	GOCACHE=$(GO_CACHE) go test ./... $(ARGS)

gui-build:
	cd packages/gui && npm run build
