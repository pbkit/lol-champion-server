# vim: set noexpandtab:
.PHONY: all clean
.DEFAULT_GOAL := all
PROTO_FILES += .pollapo/pbkit/interface-lol-champion-server/lol_champion_service.proto
SRC_FILES += src/*.go
SRC_FILES += src/**/*.go

$(PROTO_FILES): pollapo.yml
	pollapo install

src/gen: $(PROTO_FILES)
	mkdir -p src/gen
	protoc --go_out=src --go-grpc_out=src -I=.pollapo $(PROTO_FILES)

lol-champion-server: src/gen $(SRC_FILES)
	cd src; go build -o ..

all: ensure/protoc ensure/go lol-champion-server

clean:
	rm -rf .pollapo
	rm -rf src/gen
	rm -f lol-champion-server

ensure/%:
	@command -v $(*F) &>/dev/null || { echo "error: command $(*F) is not found."; exit -1; }
