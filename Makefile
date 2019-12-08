PROTOC ?= protoc

PROTO_MESSAGE := $(wildcard proto/message/*.proto)
PROTO_SERVICE := $(wildcard proto/service/*.proto)
GO_MESSAGES := $(PROTO_MESSAGE:proto/message/%.proto=message/%.pb.go)
GO_SERVICES := $(PROTO_SERVICE:proto/service/%.proto=service/%.pb.go)
CPP_MESSAGES := $(PROTO_MESSAGE:proto/message/%.proto=proto/cpp/%.pb.cc)
CPP_SERVICES := $(PROTO_SERVICE:proto/service/%.proto=proto/cpp/%.grpc.pb.cc)
CPP_MESSAGES_HEADER := $(CPP_MESSAGES:%.cc=%.h)
CPP_SERVICES_HEADER := $(CPP_SERVICES:%.cc=%.h)

# All targets needed for server
all: ${GO_MESSAGES} ${GO_SERVICES}

# Generates protobuf and grpc files for cpp
proto-cpp: proto/cpp ${CPP_SERVICES} ${CPP_MESSAGES}

# Generate protobuf golang
message/%.pb.go: proto/message/%.proto
	${PROTOC} -Iproto/message --go_out=${GOPATH}/src $<

# Generate grpc golang
service/%.pb.go: proto/service/%.proto
	${PROTOC} -Iproto -Iproto/service -Iproto/message \
		--go_out=plugins=grpc:. $<

# Create cpp folder for protobuf and grpc
proto/cpp:
	mkdir -p $@

# Generate protobuf cpp
proto/cpp/%.pb.cc: proto/message/%.proto
	${PROTOC} -Iproto/message --cpp_out=proto/cpp $<

# Generate grpc cpp
proto/cpp/%.grpc.pb.cc: proto/service/%.proto
	${PROTOC} -Iproto/service -Iproto -Iproto/message --grpc_out=proto/cpp \
		--plugin=protoc-gen-grpc=`which grpc_cpp_plugin` $<

clean:
	rm -f ${GO_MESSAGES} ${GO_SERVICES} \
		${CPP_MESSAGES} ${CPP_MESSAGES_HEADER} \
		${CPP_SERVICES} ${CPP_SERVICES_HEADER}
