package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Kethsar/getwindowprocname/proto"

	"google.golang.org/grpc"
)

func getProcessName() string {
	conn, err := grpc.Dial("192.168.1.10:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewRemoteProcClient(conn)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	proc, err := client.GetProcName(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalln(err)
	}

	return proc.GetName()
}
