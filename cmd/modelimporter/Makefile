CC = g++
CFLAGS = `pkg-config --cflags assimp protobuf`
LDFLAGS = `pkg-config --libs assimp protobuf`

all: modelimporter

protogen:
	protoc -I/usr/local/include -I../../protos --cpp_out=. ../../protos/model.proto

modelimporter: protogen
	$(CC) $(CFLAGS) $(LDFLAGS) --std=c++11 *.cc *.cpp `pkg-config --cflags assimp protobuf` `pkg-config --libs assimp protobuf` -o modelimporter

clean:
	rm modelimporter
