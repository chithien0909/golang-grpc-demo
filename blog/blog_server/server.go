package main

import (
	"../blogpb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type server struct {

}

type blogItem struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string		`bson:"author_id"`
	Content string		`bson:"content"`
	Title string		`bson:"title"`
}

func (*server) CreateBlog(_ context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	fmt.Println("Create blog request")

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Connot convert to OID"),
		)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(_ context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	fmt.Println("Read blog request")

	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse id"),
		)
	}

	data := &blogItem{}

	filter := bson.M{
		"_id": oid,
	}

	res := collection.FindOne(context.Background(), filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Not found blog with sepecial ID: %v", err),
		)
	}
	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil


}

func (*server) UpdateBlog(_ context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	oid, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse id"),
		)
	}
	data := &blogItem{}

	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Not found blog with sepecial ID: %v", err),
		)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Title = blog.GetTitle()
	data.Content = blog.GetContent()

	_, updateErr := collection.UpdateOne(context.Background(), filter, data)

	if updateErr != nil{
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Can not update object in MongoDB: %v", updateErr),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}
func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:			data.ID.Hex(),
		Title: 		data.Title,
		Content: 	data.Content,
		AuthorId: 	data.AuthorID,
	}
}
var collection *mongo.Collection

func main()  {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog service started...")

	// Connect to mongodb
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Connect to mongodb failed: %v\n", err)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received

	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	_ = lis.Close()
	fmt.Println("Closing MongoDB Connection")
	_ = client.Disconnect(ctx)
	fmt.Println("End of Program")

}


