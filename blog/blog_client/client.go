package main

import (
	"../blogpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	fmt.Println("Blog server started")

	opts := grpc.WithInsecure()

	//cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Count not connect to server: %v\n", err)

	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	//fmt.Printf("Created client: %f\n", c)
	//createBlog(c)
	readBlog(c)


}

func createBlog(c blogpb.BlogServiceClient)  {

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		Title: "My first blog",
		Content: "Content of the first blog",
		AuthorId: "Stephen",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Create blog failed: %v\n", err)
	}

	fmt.Printf("Blog has been created: %v\n", createBlogRes)
}

func readBlog(c blogpb.BlogServiceClient) {
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "215b3biu"})

	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	//blogId, _ := primitive.ObjectIDFromHex("6006cb5eadae18f2295f7401")

	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: "6006cb5eadae18f2295f7401",
	}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)

	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)
}