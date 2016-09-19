all: modelimporter

protogen:
	protoc -I/usr/local/include -I../../protos --cpp_out=. ../../protos/model.proto

modelimporter: protogen
	g++ -I/usr/local/include/assimp -I/usr/local/include -lassimp -lprotobuf -L/usr/local/lib -std=c++14 -Wall -o modelimporter *.cpp *.cc

clean:
	rm modelimporter
