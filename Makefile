binary = paper.core
GClean = go clean
GBuild = go build -v -tags=jsoniter -o $(binary)

ifeq ($(OS),Windows_NT)
 	binary = "paper.core.exe"
#else
# 	ifeq ($(shell uname),Darwin)
#  		binary = "paper.core"
# 	else
#  		binary = "paper.core"
# 	endif
endif

all: build

build:
	$(GClean)
	$(GBuild) -ldflags "-s -w"

build_auth:
	$(GClean)
	go build -v -tags=jsoniter,module_auth -o $(binary) -ldflags "-s -w"

clean:
	$(GClean)