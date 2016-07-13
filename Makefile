# Set to "-x" for verbose builds
VERBOSE ?=

static:
	@cd crane && \
		go build -o crane.static ${VERBOSE} -ldflags '-extldflags "-L/usr/lib64 -L/usr/pkg/lib -lssh2 -lssl -lcrypto -lz -lpthread -static"'

shared:
	@cd crane && go build ${VERBOSE}

clean:
	rm -fr **.core crane/crane{,.static} crane-manifest/crane-manifest **~ obj _obj
