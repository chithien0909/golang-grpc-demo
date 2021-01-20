package main

import (
	"../calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main()  {

	fmt.Println("Calculator client\n")

	cc, err := grpc.Dial("localhost:3000", grpc.WithInsecure())

	if err != nil{
		log.Fatalf("Count not connect to server: %v\n", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiStreaming(c)

	doSquareRoot(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient)  {

	fmt.Printf("Starting to do Unary RPC...\n")

	req := &calculatorpb.CalculatorRequest{
		A: 1,
		B: 5,
	}
	res, err := c.Calculator(context.Background(), req)

	if err!=nil{
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	fmt.Printf("Response from Calculator: %v\n", res.Result)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Printf("Starting to do Server Streaming RPC...\n")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 143535,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Prime Number Decomposition: %v\n", err)
	}

	log.Printf("Response from GreetManyTimes: \n")
	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error when reading stream: %v\n", err)
		}

		fmt.Printf("%d ",msg.PrimeFactor)
	}
	fmt.Printf("\n")
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Printf("Starting to do Server Streaming RPC...\n")

	stream, err := c.ComputeAverage(context.Background())

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Number: 2,
		},
		{
			Number: 2,
		},
		{
			Number: 4,
		},
		{
			Number: 2,
		},
	}


	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending number %v\n", req.GetNumber())
		err := stream.Send(req)

		if err != nil {
			log.Fatalf("Error when sending request: %v\n", err)
		}

	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.Average)
}

func doBiStreaming(c calculatorpb.CalculatorServiceClient)  {

	fmt.Printf("Starting to do BiDi Streaming RPC...\n")

	stream, err := c.FindMaximum(context.Background())


	if err != nil {
		log.Fatalf("Error while creating the stream: %v\n", err)
		return
	}

	waitc := make(chan struct{})
	numbers := []float64{1.0, 5.0, 7.0, 3, 1, 9, 4, 24}
	go func() {
		for _, num := range numbers{

			fmt.Printf("Sending a message: %v\n", num)

			errSend := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})

			if errSend != nil {
				log.Fatalf("Error while sending request to server: %v\n", errSend)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		_ = stream.CloseSend()
	}()

	go func() {
		for {
			req, errReceive := stream.Recv()
			if errReceive == io.EOF {
				break
			}


			if errReceive != nil {
				log.Fatalf("Error while receiving: %v\n", errReceive)
			}

			fmt.Printf("Max number: %v\n", req.GetMaxNumber())
		}
		close(waitc)
	}()

	<-waitc

	fmt.Printf("It's done")
}

func doSquareRoot(c calculatorpb.CalculatorServiceClient) {

	fmt.Printf("Starting to do Square RPC...\n")

	//correct call
	doErrorCall(c, 10)
	//error call
	doErrorCall(c, -2)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32)  {


	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})

	if err != nil {
		errRes, ok := status.FromError(err)

		if ok {
			// Actual error from GRPC (user error)
			fmt.Printf("Error message from the server: %q\n",errRes.Message())
			fmt.Printf(string(errRes.Code()))
			if errRes.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoor: %v\n", err)
		}
	}

	fmt.Printf("Result of SquareRoot of  %v: %v\n", n, res.GetNumberRoot())
}