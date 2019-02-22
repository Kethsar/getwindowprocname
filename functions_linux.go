package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Kethsar/getwindowprocname/proto"

	"google.golang.org/grpc"
)

func getProcessName() string {
	if len(c.ServerAddress) < 1 {
		log.Fatalln("Server address not found in config. It is needed to run the client on Linux.")
	}

	conn, err := grpc.Dial(c.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewRemoteProcClient(conn)
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	proc, err := client.GetProcName(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	return proc.GetName()
}

func startServer() {
	log.Fatalln("Server mode is for Windows only.")
}
