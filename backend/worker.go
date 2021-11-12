package main

import (
	"context"

	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/Banana-no-ana/support-ticket-tycoon/backend/clockclient"
	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	serviceutils "github.com/Banana-no-ana/support-ticket-tycoon/backend/utils"
	"github.com/gorilla/mux"
)

//Global Variables!!!!
var assignedCases []*pb.Case
var completedCases []*pb.Case
var workerID int32
var _skills *pb.WorkerSkill
var customerConn pb.CustomerClient

type Worker struct {
	Name   string //Generated from a list of names
	FaceID int    //Icon for worker face
	ID     int    //Assigned worker ID
}

func listCases(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to list all cases")

	// assignedArrary := []string{}
	// for _, v := range assignedCases {
	// 	jsonMarshal, _ := protojson.Marshal(v)
	// 	assignedArrary = append(assignedArrary, string(jsonMarshal))
	// }
	// completedArrary := []string{}
	// for _, v := range assignedCases {
	// 	jsonMarshal, _ := protojson.Marshal(v)
	// 	completedArrary = append(completedArrary, string(jsonMarshal))
	// }

	reply := &pb.ListCaseReply{AssignedCases: assignedCases, CompletedCases: completedCases}

	b, _ := protojson.Marshal(reply)
	fmt.Fprintf(w, string(b))
}

func (WorkerServer) GetCaseState(ctx context.Context, c *pb.Case) (*pb.Case, error) {
	for _, myc := range assignedCases {
		if myc.CaseID == c.CaseID {
			return myc, nil
		}
	}
	return &pb.Case{}, status.Error(codes.NotFound, "Case not found")

}

func caseAssign(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])

	c := pb.Case{CaseID: int32(caseid), Status: pb.CaseStatus_Assigned}
	assignedCases = append(assignedCases, &c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %d, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc
	return
}

func caseUnAssign(caseid int) bool {
	log.Println("Received request to unassign case: ", caseid)

	for i, c := range assignedCases {
		if c.CaseID == int32(caseid) {
			assignedCases = append(assignedCases[:i], assignedCases[i+1:]...)
			log.Println("200: Unassigned case ", caseid, " from the worker")
			return true
		}
	}
	log.Println("404: Case ", caseid, " not found on the worker")
	return false

}

func caseUnAssignRequest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])

	if caseUnAssign(caseid) {
		fmt.Fprintf(w, "200: Unassign case %d from the worker", caseid)
		return
	}
	fmt.Fprintf(w, "404: Case %d not found on the worker", caseid)
}

func connectToCustomerService() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	addr := "localhost:8006"
	conn, _ := grpc.Dial(addr, opts...)
	// defer worker_conn.Close()

	customerConn = pb.NewCustomerClient(conn)
}

func getCaseNextStage(c *pb.Case) {
	//We give the customer the current case state.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if customerConn == nil {
		connectToCustomerService()
	}
	d, _ := customerConn.CustomerReply(ctx, c)

	*c = *d //Update to the contents of the customer's reply
}

func workOnCase(c *pb.Case) {
	curStage := c.CaseStages[c.CurrentStage-1] //Stages are 1-indexed.

	if curStage.Status == pb.StageStatus_Completed {
		//Current stage is complete, let the customer know.
		c.Status = pb.CaseStatus_WOCR
	} else if curStage.Status == pb.StageStatus_Working {
		var m int = 0

		if curStage.Completedwork >= curStage.Totalwork {
			curStage.Status = pb.StageStatus_Completed
			c.Status = pb.CaseStatus_WOCR
			return
		}

		switch curStage.Type {
		case pb.SkillEnum_Troubleshoot:
			m = int(_skills.Troubleshoot)
		case pb.SkillEnum_Build:
			m = int(_skills.Build)
		case pb.SkillEnum_Tech:
			m = int(_skills.Tech)
		case pb.SkillEnum_Usage:
			m = int(_skills.Usage)
		case pb.SkillEnum_Architecture:
			m = int(_skills.Architecture)
		case pb.SkillEnum_Environment:
			m = int(_skills.Environment)
		case pb.SkillEnum_Explain:
			m = int(_skills.Explain)
		case pb.SkillEnum_Empathy:
			m = int(_skills.Empathy)
		case pb.SkillEnum_Relationship:
			m = int(_skills.Relationship)
		}

		work := 16 * math.Pow(2, float64(m-int(curStage.Difficulty)))
		curStage.Completedwork = curStage.Completedwork + int32(work)
	}
}

//When the clock ticks, it tocks us.
func tock() {
	//Step 1: Look at the current case. What state is it in?
	//Case details should be given to us on assign.
	//Case states: New, Assigned, In-Progress, Closed.

	//Work on the current case.
	if len(assignedCases) == 0 {
		return
	}

	curCase := assignedCases[0]
	switch curCase.Status {
	case pb.CaseStatus_New:
		curCase.Status = pb.CaseStatus_WOCR
	case pb.CaseStatus_WOCR:
		getCaseNextStage(curCase)
	case pb.CaseStatus_InProgress:
		workOnCase(curCase)
	case pb.CaseStatus_Closed:
		//Move on to the next case. The worker will stop caring about this case
		log.Println("Case %d is closed now", curCase.CaseID)
		assignedCases = assignedCases[1:]
		completedCases = append(completedCases, curCase)
	default:
		curCase.Status = pb.CaseStatus_New
	}
	// TODO: Update the API server

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
	_skills = sk
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) Hello(context.Context, *pb.Response) (*pb.Response, error) {
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) KillWorker(context.Context, *pb.Response) (*pb.Response, error) {
	go serviceutils.Kill2()
	return &pb.Response{Success: true}, nil
}

func (s *WorkerServer) Unassign(ctx context.Context, c *pb.Case) (*pb.Response, error) {
	if caseUnAssign(int(c.CaseID)) {
		return &pb.Response{Success: true}, nil
	}
	return &pb.Response{Success: false}, nil
}

func listskills(w http.ResponseWriter, req *http.Request) {
	m, _ := protojson.Marshal(_skills)
	fmt.Fprintf(w, string(m))
}

func main() {
	http_port := flag.String("http_port", ":9000", "set the http listneing port of the worker. ")
	rpc_port := flag.String("rpc_port", ":10000", "set the rpc listneing port of the worker. ")
	worker_id_flag := flag.Int("worker_id", 0, "Identify the ID of the worker. ")

	flag.Parse()

	workerID = int32(*worker_id_flag)

	_skills = &pb.WorkerSkill{}

	r := mux.NewRouter()
	r.HandleFunc("/case/assign/{caseid}", caseAssign)
	r.HandleFunc("/case/unassign/{caseid}", caseUnAssignRequest)
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
