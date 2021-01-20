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
	"time"
)

func main()  {
	fmt.Println("Hello I'm a client")

	tls := false

	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate

		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")

		if  sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	//cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Count not connect to server: %v\n", err)

	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f\n", c)

	doUnary(c)

	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiStreaming(c)
	//doUnaryWithDeadline(c, 5 * time.Second)
	//doUnaryWithDeadline(c, 1 * time.Second)
}

func doUnary(c greetpb.GreetServiceClient)  {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient)  {
	fmt.Printf("Starting to do Server Streaming RPC...\n")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName: "Maarek",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err!=nil{
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we are reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Errror while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}



}

func doClientStreaming(c greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 1",
				LastName: "Lastname 1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 2",
				LastName: "Lastname 2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 3",
				LastName: "Lastname 3",
			},
		},
	}

	fmt.Printf("Starting to do Client Streaming RPC...\n")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	for _, req := range requests{
		fmt.Printf("Sending req: %v\n", req)
		errSend := stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
		if errSend != nil {
			log.Fatalf("Error while sending request to LongGreet: %v", errSend)
		}
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiStreaming(c greetpb.GreetServiceClient) {

	fmt.Printf("Starting to do BiDi Streaming RPC...\n")

	// We create a stream by invoked the client

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating the stream: %v\n", err)
		return
	}
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 1",
				LastName:  "Lastname 1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 2",
				LastName:  "Lastname 2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 3",
				LastName:  "Lastname 3",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Firstname 4",
				LastName:  "Lastname 4",
			},
		},
	}
	// wait channel
	waitc := make(chan struct{})

	// We send a bunch of messages to the client (go routine)

	go func() {
		// function to send a bunch of messages
		for _, req := range requests {

			fmt.Printf("Sending a message: %v\n", req)

			errSend := stream.Send(req)
			if errSend != nil {
				log.Fatalf("Error while sending request to server: %v\n", errSend)
			}
			time.Sleep(1000 * time.Millisecond)

		}
		_ = stream.CloseSend()

	}()

	// We receive a bunch of messages from the client (go routine)

	go func() {
		// function to receive a bunch of messages
		for {

			res, errReceive := stream.Recv()
			if errReceive == io.EOF {

				close(waitc)
				break
			}
			if errReceive != nil {
				log.Fatalf("Error while receiving: %v\n", errReceive)
			}
			fmt.Printf("Receiving %v\n", res.GetResult())
		}
	}()

	// Block until everything is done
	<-waitc

	fmt.Println("It done")
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration)  {
	fmt.Printf("Starting to do Unary RPC...\n")
	// From greet.pb.go
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName: "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err!=nil{

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("Timeout was hit! Deadline was exceeded: %v\n", statusErr.Message())
			} else {
				fmt.Printf("Unexpected error: %v\n", err)
			}
		} else {

			log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	fmt.Printf("Response from GreetWithDeadline: %v\n", res.Result)
}