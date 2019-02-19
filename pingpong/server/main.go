package main

import "context"
import "fmt"
import "google.golang.org/grpc"
import "net"
import "flag"
import "github.com/t-drk/grpc_performance/pingpong"
import "log"

var (
	port = flag.Int("port", 50001, "The server port to listen on.")
)

type pingPong struct{}

func (p *pingPong) SayHello(ctx context.Context, request *pingpong.Hello) (*pingpong.Hello, error) {
	response := &pingpong.Hello{Message: "Reply: " + request.Message}
	return response, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("net.Listen(_) = _, %v", err)
	}
	grpcServer := grpc.NewServer()
	pingpong.RegisterPingPongServer(grpcServer, &pingPong{})
	log.Fatalln(grpcServer.Serve(lis))
}
