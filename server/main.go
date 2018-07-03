package main

import (
	"fmt"
	"net/http"

	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/takuoki/grpc-gateway-sample/proto"
)

type sampleService struct {
}

func (*sampleService) GetSample(ctx context.Context, s *pb.Sample) (*pb.Sample, error) {
	fmt.Printf("-> %+v\n", s)
	return s, nil
}

func main() {

	httpPort := ":8080"
	grpcPort := ":9090"

	// http server
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterSampleServiceHandlerFromEndpoint(context.Background(), mux, grpcPort, opts)
	if err != nil {
		log.Fatalf("register endpoint failed: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("failed to listen and serve http: %v", err)
		}
	}()

	// grpc server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterSampleServiceServer(server, new(sampleService))
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
