#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include "pingpong.grpc.pb.h"
#include <chrono>
#include <cstdlib>
#include <unistd.h>

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using pingpong::Hello;
using pingpong::PingPong;
using std::chrono::duration_cast;
using std::chrono::high_resolution_clock;
using std::chrono::microseconds;


int main(int argc, char **argv){
	std::string message("hello");
	std::string serverAddr("localhost:50001");
	int sampleCount = 1000;
	int c ;
	while ( (c = getopt(argc, argv, "s:m:c:")) != - 1) {
		switch (c) {
			case 's':
				// Server address
				serverAddr = std::string(optarg);
				break;
			case 'm':
				message = std::string(optarg);
				break;
			case 'c':
				sampleCount = std::atoi(optarg);
				break;
		}
	}

	std::unique_ptr<PingPong::Stub> client(PingPong::NewStub(grpc::CreateChannel(serverAddr, grpc::InsecureChannelCredentials())));
	for (int i = 0 ;i < sampleCount; i ++){
		Hello request;
		request.set_message(message);
		ClientContext context;
		Hello reply;
		auto start = high_resolution_clock::now();
		Status status = client->SayHello(&context, request, &reply);
		if (!status.ok()){
			std::cout << status.error_message() << std::endl;
		}
		auto stop = high_resolution_clock::now();
		auto duration = duration_cast<microseconds>(stop-start);
		std::cout << duration.count() << std::endl ;
	} 
	return 0;
}
