package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/Kethsar/getwindowprocname/proto"

	"google.golang.org/grpc"
)

// Make a remote call to a Windows machine to get the process name for the window currently below the cursor
func getWindowInfo(x, y int) *pb.WindowInfo {
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

	winfo, err := client.GetWindowInfo(ctx, &pb.Cursor{X: int32(x), Y: int32(y)})
	if err != nil {
		log.Fatalln(err)
	}

	if c.WriteToFile {
		writeWinInfoToFile(winfo)
	}

	return winfo
}

func writeWinInfoToFile(winfo *pb.WindowInfo) {
	jsbytes, err := json.MarshalIndent(winfo, "", "\t")
	if err != nil {
		log.Println("Error converting data to JSON:", err)
		return
	}

	err = ioutil.WriteFile(c.FileLocation, jsbytes, 0666)
	if err != nil {
		log.Println("Error writing json to tmp file:", err)
	}
}

func startServer() {
	log.Fatalln("Server mode is for Windows only.")
}
