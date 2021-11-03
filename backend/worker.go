package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/Banana-no-ana/support-ticket-tycoon/backend/clockclient"
	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	serviceutils "github.com/Banana-no-ana/support-ticket-tycoon/backend/utils"
	"github.com/gorilla/mux"
)

//Global Variables!!!!
var assignedCases []*pb.Case
var workerID int32
var _skills pb.WorkerSkill

type Worker struct {
	Name   string //Generated from a list of names
	FaceID int    //Icon for worker face
	ID     int    //Assigned worker ID
}

func listCases(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to list all cases")

	caseArrary := []pb.Case{}
	for _, v := range assignedCases {
		caseArrary = append(caseArrary, *v)
	}
	b, _ := json.Marshal(caseArrary)
	fmt.Fprintf(w, string(b))
}

func caseAssign(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])

	c := pb.Case{CaseID: int32(caseid), Status: "Assigned"}
	assignedCases = append(assignedCases, &c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %d, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc
	return
}

func unassign(w http.ResponseWriter, req *http.Request) {
	//TODO: Implement this
	return
}

func getCaseNextStage(c *pb.Case) {
	//We give the customer the current case state.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d, _ := customerConn.CustomerReply(ctx, c)

	d.Status = "In-Progress"
	*c = *d //Update to the contents of the customer's reply
}

func calcStageProgress(stage *pb.CaseStage) *pb.CaseStage {

}

func workOnCase(c *pb.Case) bool {
	curStage := c.CaseStages[c.CurrentStage-1]
	if curStage.Status == "Complete" {
		//Current stage is complete, let the customer know.
		c.Status = "Waiting for Customer Response"
	} else if curStage.Status == "In-Progress" {
		calcStageProgress(curStage)
	}

	return false

}

//When the clock ticks, it tocks us.
func tock() {
	//Step 1: Look at the current case. What state is it in?
	//Case details should be given to us on assign.
	//Case states: New, Assigned, In-Progress, Closed.

	//Work on the current case.
	curCase := assignedCases[0]
	switch curCase.Status {
	case "New":
		curCase.Status = "Waiting for Customer Reply"
	case "Waiting for Customer Reply":
		getCaseNextStage(curCase)
	case "In-Progress":
		workOnCase(curCase)
	case "Closed":
		assignedCases = assignedCases[1:]
		//Spend cycle to move to the next case.
	}

	// Update the API server

	return
}

type WorkerServer struct {
	pb.UnimplementedWorkerServer
}

func newWorkerServer() *WorkerServer {
	s := &WorkerServer{}
	return s
}

func (s *WorkerServer) Assign(ctx context.Context, in *pb.Case) (*pb.Response, error) {
	in.Assignee = workerID
	assignedCases = append(assignedCases, in)
	log.Println("Case ", in.GetCaseID(), "has been assigned")
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) SetWorkerSkills(c context.Context, sk *pb.WorkerSkill) (*pb.Response, error) {
	_skills = pb.WorkerSkill(*sk)
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) Hello(context.Context, *pb.Response) (*pb.Response, error) {
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) KillWorker(context.Context, *pb.Response) (*pb.Response, error) {
	go serviceutils.Kill2()
	return &pb.Response{Success: true}, nil
}

func listskills(w http.ResponseWriter, req *http.Request) {
	m, _ := json.Marshal(_skills)
	fmt.Fprintf(w, string(m))
}

func main() {
	http_port := flag.String("http_port", ":9000", "set the http listneing port of the worker. ")
	rpc_port := flag.String("rpc_port", ":10000", "set the rpc listneing port of the worker. ")
	worker_id_flag := flag.Int("worker_id", 0, "Identify the ID of the worker. ")

	flag.Parse()

	workerID = int32(*worker_id_flag)

	r := mux.NewRouter()
	r.HandleFunc("/assign/{caseid}", caseAssign)
	r.HandleFunc("/healthz", serviceutils.Healthz)
	r.HandleFunc("/kill", serviceutils.Kill)
	r.HandleFunc("/skill/list", listskills)
	r.HandleFunc("/case/list", listCases)
	// r.HandleFunc("/case/{caseid}", getCase) TODO: Implement this
	// r.HandleFunc("/unassign/{caseid}", caseAssign)

	http.Handle("/", r)
	log.Println("listening for http requests on :", *http_port)
	go http.ListenAndServe(*http_port, nil)

	go clockclient.CreateClockClient(tock)

	lis, err := net.Listen("tcp", *rpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, newWorkerServer())
	log.Printf("Worker %d serving rpc requests on: %v", workerID, lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// api_conn, _ := grpc.Dial("localhost:8002")
	// def api_conn.Close()

}
