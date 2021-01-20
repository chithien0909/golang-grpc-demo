package main

import (
	"../blogpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
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
	//readBlog(c)
	//updateBlog(c)
	//deleteBlog(c)
	listBlog(c)


}

func createBlog(c blogpb.BlogServiceClient)  {

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		Title: "My third blog",
		Content: "Content of the third blog",
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

func updateBlog(c blogpb.BlogServiceClient) {

	fmt.Println("Update the blog")
	newBlog := &blogpb.Blog{
		Id: "6007fc51a21a2e76be6cfc71",
		Title: "My First blog (edited)",
		Content: "Content of the third blog",
		AuthorId: "Stephen",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})

	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", updateErr)
	}

	fmt.Printf("Blog was update: %v\n", updateRes)

}

func deleteBlog(c blogpb.BlogServiceClient) {
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: "6007fc66a21a2e76be6cfc72"})

	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
		return
	}


	fmt.Printf("Blog id was delete: %v\n", res)
}


func listBlog(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatalf("Error while calling BlogList: %v\n", err)

	}
	for {
		listBlogRes, listBlogErr := stream.Recv()
		if listBlogErr == io.EOF {
			break
		}
		if listBlogErr != nil {

			log.Fatalf("Error while receiving list blog: %v\n", err)
		}
		fmt.Printf("Blog: %v\n", listBlogRes.GetBlog())
	}
}

