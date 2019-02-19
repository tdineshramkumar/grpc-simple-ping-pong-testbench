package main

import "context"
import "fmt"
import "github.com/t-drk/grpc_performance/pingpong"
import "google.golang.org/grpc"
import "flag"
import "log"
import "time"

var (
	serverAddr  = flag.String("s", "127.0.0.1:50001", "The server address of the form host:port")
	sampleCount = flag.Int("c", 1000, "The number of pings to server")
	message     = flag.String("m", "hello", "the message to send to server")
)

func StartTimer() func() {
	/*
		Prints time in microseconds
	*/
	start := time.Now()
	return func() {
		d := time.Now().Sub(start)
		fmt.Printf("%d\n", d.Nanoseconds()/1000)
	}
}

func SayHello(client pingpong.PingPongClient, message *pingpong.Hello) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stopfunc := StartTimer()
	if _, err := client.SayHello(ctx, message); err != nil {
		log.Fatalf("client.SayHello(_) = _, %v", err)
	}
	stopfunc()
}
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial() = _, %v", err)
	}
	defer conn.Close()
	client := pingpong.NewPingPongClient(conn)
	message := &pingpong.Hello{Message: *message}
	for i := 0; i < *sampleCount; i++ {
		SayHello(client, message)
	}
}
