package main

import (
	"../calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct{}


func (*server) Calculator(_ context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator function was invoked with %v\n", req)
	// a := req.GetCalculatorArgs().A
	// b := req.GetCalculatorArgs().B
	a := req.A
	b := req.B

	result := a + b

	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	number := req.GetNumber()

	k := int64(2)

	for number > 1 {

		if number % k == 0 {
			err := stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: k,
			})
			if err != nil{
				fmt.Printf("Sending response to client error: %v", err);
			}
			number = number / k
		} else {
			k++
			fmt.Printf("Devision has increased to %v\n", k)
		}

	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked\n")

	var sum float64
	var count float64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&calculatorpb.ComputeAverageResponse{
					Average: sum / count,
			})
		}
		if err != nil {
			log.Fatalf("Error while recevied request from ComputeAverage: %v\n", err)
		}
		count += 1
		sum += float64(req.GetNumber())
	}

}


func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	max := float64(-9999)
	for {
		req, err := stream.Recv()


		if err == io.EOF {
			return nil
		}

		if err != nil {

			log.Fatalf("Error while reading client stream: %v", err)
		}
		num := req.GetNumber()

		if num > max {
			max = num
			errSend := stream.Send(&calculatorpb.FindMaximumResponse{
				MaxNumber: max,
			})

			if errSend != nil {
				log.Fatalf("Error while sending response: %v\n", errSend)
			}
		}

	}
}


func (*server) SquareRoot(_ context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v\n", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "localhost:3000")

	if err != nil {

		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	//Register reflection service on gRPC server
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {

		log.Fatalf("Failed to serve: %v", err)
	}

}
