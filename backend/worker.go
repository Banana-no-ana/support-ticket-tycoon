package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/Banana-no-ana/support-ticket-tycoon/backend/clockclient"
	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	"github.com/gorilla/mux"
)

//Global Variables!!!!
var assignedCases []Case
var workerID int
var _skills pb.WorkerSkill

type Case struct {
	CaseID int32
}

type Worker struct {
	Name   string //Generated from a list of names
	FaceID int    //Icon for worker face
	ID     int    //Assigned worker ID
}

func caseAssign(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])

	c := Case{CaseID: int32(caseid)}
	assignedCases = append(assignedCases, c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %d, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc
	return
}

func unassign(w http.ResponseWriter, req *http.Request) {
	//TODO: Implement this
	return
}

//When the clock ticks, it tocks us.
func tock() {
	log.Println("We've been Tock'd ")
}

type WorkerServer struct {
	pb.UnimplementedWorkerServer
}

func newWorkerServer() *WorkerServer {
	s := &WorkerServer{}
	return s
}

func (s *WorkerServer) Assign(ctx context.Context, in *pb.Case) (*pb.Response, error) {
	c := Case{CaseID: in.GetCaseID()}
	assignedCases = append(assignedCases, c)
	log.Println("Case : ", in.GetCaseID(), "was assigned")
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) SetWorkerSkills(c context.Context, sk *pb.WorkerSkill) (*pb.Response, error) {
	_skills = pb.WorkerSkill(*sk)
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) Hello(context.Context, *pb.Response) (*pb.Response, error) {
	return &pb.Response{Success: true}, nil
}

func main() {
	http_port := flag.String("http_port", ":9000", "set the listneing port of the worker. ")
	rpc_port := flag.String("rpc_port", ":10000", "set the listneing port of the worker. ")
	worker_id_flag := flag.Int("worker_id", 1, "Identify the ID of the worker. ")

	flag.Parse()

	workerID = *worker_id_flag

	r := mux.NewRouter()
	r.HandleFunc("/assign/{caseid}", caseAssign)
	// r.HandleFunc("/unassign/{caseid}", caseAssign)

	http.Handle("/", r)
	log.Println("listening on :", *http_port)
	go http.ListenAndServe(*http_port, nil)

	go clockclient.CreateClockClient(tock)

	lis, err := net.Listen("tcp", *rpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, newWorkerServer())
	log.Printf("Worker %d started at %v", workerID, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// api_conn, _ := grpc.Dial("localhost:8002")
	// def api_conn.Close()

}
