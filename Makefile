PROTOC ?= /usr/local/bin/protoc

PROTO_MESSAGE := $(wildcard internal/proto/message/*.proto)
PROTO_SERVICE := $(wildcard internal/proto/service/*.proto)
GO_MESSAGES := $(PROTO_MESSAGE:internal/proto/message/%.proto=internal/message/%.pb.go)
GO_SERVICES := $(PROTO_SERVICE:internal/proto/service/%.proto=internal/grpc/service/%.pb.go)
CPP_MESSAGES := $(PROTO_MESSAGE:internal/proto/message/%.proto=internal/proto/cpp/%.pb.cc)
CPP_SERVICES := $(PROTO_SERVICE:internal/proto/service/%.proto=internal/proto/cpp/%.grpc.pb.cc)
CPP_MESSAGES_HEADER := $(CPP_MESSAGES:%.cc=%.h)
CPP_SERVICES_HEADER := $(CPP_SERVICES:%.cc=%.h)

# All targets needed for server
all: ${GO_MESSAGES} ${GO_SERVICES}

# Generates protobuf and grpc files for cpp
proto-cpp: internal/proto/cpp ${CPP_SERVICES} ${CPP_MESSAGES}

# Generate protobuf golang
internal/message/%.pb.go: internal/proto/message/%.proto
	${PROTOC} -Iinternal/proto/message --go_out=${GOPATH}/src $<

# Generate grpc golang
internal/grpc/service/%.pb.go: internal/proto/service/%.proto
	${PROTOC} -Iinternal/proto -Iinternal/proto/service -Iinternal/proto/message \
		--go_out=plugins=grpc:internal/grpc $<

# Create cpp folder for protobuf and grpc
internal/proto/cpp:
	mkdir -p $@

# Generate protobuf cpp
internal/proto/cpp/%.pb.cc: internal/proto/message/%.proto
	${PROTOC} -Iinternal/proto/message --cpp_out=internal/proto/cpp $<

# Generate grpc cpp
internal/proto/cpp/%.grpc.pb.cc: internal/proto/service/%.proto
	${PROTOC} -Iinternal/proto/service -Iinternal/proto -Iinternal/proto/message --grpc_out=internal/proto/cpp \
		--plugin=protoc-gen-grpc=`which grpc_cpp_plugin` $<

clean:
	rm -f ${GO_MESSAGES} ${GO_SERVICES} ${CPP_MESSAGES} ${CPP_MESSAGES_HEADER} ${CPP_SERVICES} ${CPP_SERVICES_HEADER}
