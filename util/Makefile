all: build

static:
	@go build -x -ldflags '-extldflags "-lpthread -lz -lssl -lcrypto -lssh2 -static"'

build:
	@go build

fmt:
	@gofmt -w *.go

clean:
	rm -fr **.core crane/crane crane-manifest/crane-manifest **~ obj _obj
