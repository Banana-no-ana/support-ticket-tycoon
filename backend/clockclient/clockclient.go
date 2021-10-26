package clockclient

import (
	"context"
	"io"
	"log"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	"google.golang.org/grpc"
)

func CreateClockClient(tock func()) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	clock_conn, _ := grpc.Dial("localhost:8000", opts...)
	defer clock_conn.Close()
	clock_client := pb.NewClockClient(clock_conn)
	registerwithClock(clock_client, tock)
}

func registerwithClock(client pb.ClockClient, tock func()) {
	log.Println("Connecting to the clock server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.Register(ctx, &pb.WorkerRegister{ID: "worker-1"})
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		tock()
	}
}
