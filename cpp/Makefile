CC = g++ --std=c++11
CFLAGS = -O3 `pkg-config --libs grpc protobuf grpc++`
DEPS = pingpong.grpc.pb.h pingpong.pb.h

all: client server


pingpong.grpc.pb.o:	pingpong.grpc.pb.cc $(DEPS)
	$(CC) -c  $(DEPS) $< $(CFLAGS)

pingpong.pb.o:	pingpong.pb.cc $(DEPS)
	$(CC) -c  $(DEPS) $< $(CFLAGS)

client:	client.cc pingpong.grpc.pb.o pingpong.pb.o
	$(CC) -o client client.cc pingpong.grpc.pb.o pingpong.pb.o $(CFLAGS)

server:	server.cc pingpong.grpc.pb.o pingpong.pb.o
	$(CC) -o server server.cc pingpong.grpc.pb.o pingpong.pb.o $(CFLAGS)

async_server:	async_server.cc pingpong.grpc.pb.o pingpong.pb.o
	$(CC) -o async_server async_server.cc pingpong.grpc.pb.o pingpong.pb.o $(CFLAGS)
clean:	
	rm *.o server client

