PROJECT = sample

BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
UNAME := $(shell uname -s)

ifeq ($(UNAME),Darwin)
#	PKG_CONFIG_PATH = "/usr/local/opt/readline/lib/pkgconfig"
	LDFLAGS="-L./lib/macos -rpath ."
	CFLAGS="-I../include"
	DYLIB="./lib/macos/libtcod.dylib"
endif


GO_FILES = $(shell find . ../ -type f -name '*.go')
GO_TEST_FLAGS ?= -p 8 -count=1 -cover

VERBOSITY ?= 2

default: $(PROJECT)

$(PROJECT): $(GO_FILES)
	@CGO_LDFLAGS=$(LDFLAGS) CGO_CFLAGS=$(CFLAGS) go build -o ./$(PROJECT)

.PHONY: dist
dist: $(PROJECT)
	@mkdir -p dist
	@cp $(PROJECT) dist
	@cp $(DYLIB) dist
	@cp -r data dist/data

# Launch a server with TLS support
.PHONY: run
run: $(PROJECT)
	@cd dist && ./$(PROJECT) --font=./data/fonts/terminal16x16_gs_ro.png --font.rows

.PHONY: clean
clean::
	@rm $(PROJECT)
	@rm -rf dist


