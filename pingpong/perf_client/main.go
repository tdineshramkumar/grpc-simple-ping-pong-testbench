package main

import "context"
import "fmt"
import "github.com/t-drk/grpc_performance/pingpong"
import "google.golang.org/grpc"
import "flag"
import "log"
import "time"

var (
	serverAddr  = flag.String("server_addr", "127.0.0.1:50001", "The server address of the form host:port")
	message     = flag.String("message", "hello", "the message to send to server")
	numRoutines = flag.Int("num_routines", 1, "the number of concurrent go routines to send requests")
	duration    = flag.Float64("duration", 60.0, "The duration of the test")
)

func init() {
	flag.Parse()
	fmt.Printf("Sending message: [%v] to server: [%v] for duration: [%v] using [%v] routines\n", *message, *serverAddr, *duration, *numRoutines)

}

const RPCTimeout = 10 * time.Second

type TimerFunc func() time.Duration

type Metrics struct {
	TotalDuration float64
	NumRequests   int
	NumErrors     int
}

func Timer() TimerFunc {
	start := time.Now()
	return TimerFunc(func() time.Duration {
		return time.Now().Sub(start)
	})
}

func SayHello(client pingpong.PingPongClient) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()
	message := &pingpong.Hello{Message: *message}
	timerFunc := Timer()
	_, err := client.SayHello(ctx, message)
	return timerFunc(), err
}

type TaskFunc func() (time.Duration, error)

func SingleLoadSession(fn TaskFunc, metricsChan chan<- Metrics) {
	timerFunc := Timer()
	var metrics Metrics
	for timerFunc().Seconds() <= float64(*duration) {
		d, err := fn()
		if err != nil {
			metrics.NumErrors++
			continue
		}
		metrics.TotalDuration += d.Seconds()
		metrics.NumRequests++
	}
	metricsChan <- metrics
}

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial() = _, %v", err)
	}
	defer conn.Close()
	client := pingpong.NewPingPongClient(conn)

	metricsChan := make(chan Metrics, *numRoutines)
	taskFunc := TaskFunc(func() (time.Duration, error) {
		return SayHello(client)
	})
	for i := 0; i < *numRoutines; i++ {
		go SingleLoadSession(taskFunc, metricsChan)
	}
	var globalMetrics Metrics
	// Obtain the metrics from each of the routines
	for i := 0; i < *numRoutines; i++ {
		metrics := <-metricsChan
		globalMetrics.NumErrors += metrics.NumErrors
		globalMetrics.NumRequests += metrics.NumRequests
		globalMetrics.TotalDuration += metrics.TotalDuration
	}
	// Average Duration
	avgDuration := globalMetrics.TotalDuration / float64(*numRoutines)
	avgRequestRate := float64(globalMetrics.NumRequests) / avgDuration
	avgRequestTime := globalMetrics.TotalDuration / float64(globalMetrics.NumRequests)
	fmt.Println("Avg. Request Rate: ", avgRequestRate, "requests/second")
	fmt.Println("Avg. Request Time: ", avgRequestTime*1e6, "microseconds")
	fmt.Println("Num. Errors: ", globalMetrics.NumErrors)
}
