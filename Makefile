PROJECT = pg_stats
VERSION ?= 0.0.1

REPO = github.com/sbowman/pgstats
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
UNAME := $(shell uname -s)

ifeq ($(UNAME),Darwin)
	PKG_CONFIG_PATH = "/usr/local/opt/readline/lib/pkgconfig"
endif

GO_FILES = $(shell find . -path ./.idea -prune -o -type f -name '*.go' -print)

default: $(PROJECT)

$(PROJECT): $(GO_FILES)
	@cd cli && go build -o ../$(PROJECT)
