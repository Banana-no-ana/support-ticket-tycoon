package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	"github.com/gorilla/mux"
)

//Global Variable!!!!
var assignedCases []Case
var workerID int

type Case struct {
	CaseID   int32
	State    string //State is the current case state
	Assignee int
}

type Worker struct {
	Name   string //Generated from a list of names
	FaceID int    //Icon for worker face
	ID     int    //Assigned worker ID
}

func calcNextTick() {
	//Take the first assigned case, and figure out what to do with it.
	if len(assignedCases) == 0 {
		return
	}
	currentCase := assignedCases[0]

	switch currentCase.State {
	case "Assigned": //Just assigned. need to triage.
		return
	}
	return
}

// func caseAssign(w http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	caseid, _ := strconv.Atoi(vars["caseid"])

// 	c := Case{CaseID: caseid, State: "Assigned"}
// 	assignedCases = append(assignedCases, c)
// 	log.Println("case assigned: ", caseid)
// 	fmt.Fprintf(w, "case accepted %d, assigned cases: %d \n", caseid, len(assignedCases))

// 	//Trigger next work item recalc
// 	return
// }

// work on the next assigned thing.
func tick(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tickNum := vars["ticknum"]
	log.Println("We've been Ticked! Ticknum: ", tickNum)

	//TODO: do the actual actions of the tick.
	calcNextTick()
	return
}

func unassign(w http.ResponseWriter, req *http.Request) {
	//TODO: Implement this
	return
}

func registerwithClock(client pb.ClockClient) {
	log.Println("Connecting to the clock server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.Register(ctx, &pb.WorkerRegister{ID: "worker-1"})
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		TickTick, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Printf("Tick received: %d", TickTick.TickNum)
	}
}

type WorkerServer struct {
	pb.UnimplementedWorkerServer
}

func newServer() *WorkerServer {
	s := &WorkerServer{}
	return s
}

func (s *WorkerServer) Assign(ctx context.Context, in *pb.Case) (*pb.Response, error) {
	c := Case{CaseID: in.GetCaseID(), State: "Assigned", Assignee: workerID}
	assignedCases = append(assignedCases, c)
	log.Println("Case : ", in.GetCaseID(), "was assigned")
	return &pb.Response{Success: true}, nil
}

func main() {
	http_port := flag.String("http_port", ":8081", "set the listneing port of the worker. ")
	rpc_port := flag.String("rpc_port", ":8080", "set the listneing port of the worker. ")
	worker_id_flag := flag.Int("worker_id", 1, "Identify the ID of the worker. ")

	flag.Parse()

	workerID = *worker_id_flag

	r := mux.NewRouter()
	// r.HandleFunc("/assign/{caseid}", caseAssign)
	// r.HandleFunc("/unassign/{caseid}", caseAssign)
	r.HandleFunc("/tick/{ticknum}", tick)

	http.Handle("/", r)
	log.Println("listening on :", *http_port)
	go http.ListenAndServe(*http_port, nil)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	clock_conn, _ := grpc.Dial("localhost:8000", opts...)
	defer clock_conn.Close()
	clock_client := pb.NewClockClient(clock_conn)
	go registerwithClock(clock_client)

	lis, err := net.Listen("tcp", *rpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, newServer())
	log.Printf("Worker %d started at %v", workerID, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// api_conn, _ := grpc.Dial("localhost:8002")
	// def api_conn.Close()

}
