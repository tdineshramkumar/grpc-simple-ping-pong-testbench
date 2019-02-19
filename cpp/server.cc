
#include <iostream>
#include <string>
#include <memory>
#include <grpcpp/grpcpp.h>
#include "pingpong.grpc.pb.h"
using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using pingpong::PingPong;
using pingpong::Hello;

class PingPongImpl final : public PingPong::Service{
	Status SayHello(ServerContext* context, const Hello *request, Hello *response) override {
		std::string prefix("Reply: ");
		response->set_message(prefix + request->message());
		return ::grpc::Status::OK;
	}
};

int main(int argc, char **argv) {
	std::string server_address("0.0.0.0:50001");
	PingPongImpl service;
	ServerBuilder builder;
	builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
	builder.RegisterService(&service);
	std::unique_ptr<Server> server(builder.BuildAndStart());
	server->Wait();
	return 0;
}
