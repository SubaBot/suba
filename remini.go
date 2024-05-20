package main

import (
	"context"
	"time"

	pb "github.com/SubaBot/suba/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Toxicity allows to check if a string contains
// toxicity (such as insult).
func Toxicity(text string) (*pb.ToxicityReply, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &pb.ToxicityReply{}, err
	}

	defer conn.Close()
	c := pb.NewReminiClient(conn)

	// Contact the server and print out its response
	// If no response in 2 seconds, cancel it.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	// Make request.
	r, err := c.Toxicity(ctx, &pb.StringRequest{
		String_: text,
	})

	if err != nil {
		return &pb.ToxicityReply{}, err
	}

	return r, nil
}
