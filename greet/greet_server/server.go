package main

import (
	"../greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
}


func (*server) Greet(_ context.Context,req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}


func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().FirstName
	for i := 0; i<10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)

		if err != nil{
			log.Fatalf("Send stream data error: %v", err)
		}

		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked")
	var result string
	for {

		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {

			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()

		result += "Hello " + firstName + "! "

	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {

			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if sendErr != nil {

			log.Fatalf("Error while sending response to client stream: %v", err)
		}
	}
	return nil
}

func (*server) GreetWithDeadline(ctx context.Context, req * greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)

	for i:=0; i<3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()

	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main()  {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	tls := false

	if tls {

		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

		if sslErr != nil{
			log.Fatalf("Failed loading certificates: %v\n", sslErr)
			return
		}

		opts = append(opts, grpc.Creds(creds))
	}


	s := grpc.NewServer(opts...)

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}