//Customers respond to case messages and may assign csat to close case.
//Customers do not create cases. That's done by the case generator
//API polls the customer on their sentiment.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/Banana-no-ana/support-ticket-tycoon/backend/clockclient"
	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	serviceutils "github.com/Banana-no-ana/support-ticket-tycoon/backend/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var _difficulty *pb.Difficulty
var cases map[int32]*pb.Case //Mapping caseIDs to case objects.

//When the clock ticks, it tocks us.
func tock() {
	//Calculate if we need to mark a case as failed.
	// log.Println("We've been tocked!")
}

type CustomerServer struct {
	pb.UnimplementedCustomerServer
}

func newCustomerServer() *CustomerServer {
	s := &CustomerServer{}
	return s
}

func (CustomerServer) SetDifficulty(ctx context.Context, dif *pb.Difficulty) (*pb.Response, error) {
	_difficulty = dif
	return &pb.Response{Success: true}, nil
}
func (CustomerServer) RegisterCase(ctx context.Context, c *pb.Case) (*pb.Response, error) {
	log.Println("Received request to register case: " + strconv.Itoa(int(c.CaseID)))
	if cases[c.CaseID] == nil {
		cases[c.CaseID] = c
		return &pb.Response{Success: true}, nil
	}

	return &pb.Response{Success: false}, errors.New("Case " + string(c.CaseID) + " already registered")
}

func getDifficulty() int32 {

	var dif int32 = 0

	if _difficulty == nil {
		_difficulty = &pb.Difficulty{}
	}

	if _difficulty.MaxDifficulty == _difficulty.MinDifficulty {
		dif = _difficulty.MaxDifficulty
	} else {
		dif = rand.Int31n(_difficulty.MaxDifficulty-_difficulty.MinDifficulty) + _difficulty.MinDifficulty
	}

	return dif

}

func generateCaseStage(stage int32, last_stage bool) *pb.CaseStage {
	dif := getDifficulty()

	var typ int = 0

	if stage == 1 {
		//Initial stage is only the first 3 types.
		typ = rand.Intn(3)
		dif = 0 //Initial triage is always short.
	} else if last_stage {
		//Only use the last 3 for final stages
		typ = rand.Intn(3) + 6
	} else {
		typ = rand.Intn(9)
	}

	st := pb.CaseStage{StageID: stage, Status: pb.StageStatus_Working, Difficulty: dif,
		Totalwork: 192 + dif*64, Completedwork: 0, Type: pb.SkillEnum(typ)}
	return &st
}

func (CustomerServer) CustomerReply(ctx context.Context, c *pb.Case) (*pb.Case, error) {
	dif := getDifficulty()

	//We calculate the number of stages here.

	//First stage always one of the first 3.
	//Last stage always one of the last 3.
	//Middle stages can be any of them.
	if c.Status == pb.CaseStatus_WOCR {
		//Need to have at least 3 stages, and the same number of stages as the "difficulty"
		if c.CurrentStage < 3 || c.CurrentStage < dif {
			c.CurrentStage = c.CurrentStage + 1

			//This may end up with more than one "last stage", which is fine.
			log.Println("Generating stage ", strconv.Itoa(int(c.CurrentStage)),
				" for case: ", strconv.Itoa(int(c.CaseID)))
			curStage := generateCaseStage(c.CurrentStage,
				c.CurrentStage >= dif && c.CurrentStage >= 3)
			c.CaseStages = append(c.CaseStages, curStage)
			c.Status = pb.CaseStatus_InProgress
			return c, nil
		} else {
			//currentstage >= difficulty. That was the last stage, close the case.
			c.Status = pb.CaseStatus_Closed
			return c, nil
		}

	}
	return c, errors.New("Customer not ready to reply to a case that's not waiting for their reply")
}

func listCases(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to list all cases")

	caseArrary := []pb.Case{}
	for _, v := range cases {
		caseArrary = append(caseArrary, *v)
	}
	b, _ := json.Marshal(caseArrary)
	fmt.Fprintf(w, string(b))
}

//HTTP request to manaully register a case so we know about it.
func registerCase(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to manually register a case")

	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])
	caseid32 := int32(caseid)

	cases[caseid32] = &pb.Case{CaseID: caseid32}
	fmt.Fprintf(w, "Case %d register successful", caseid32)
}

func main() {
	cases = make(map[int32]*pb.Case)
	http_port := flag.String("http_port", ":8005", "set the http listneing port of the customer server. ")
	rpc_port := flag.String("rpc_port", ":8006", "set the rpc listneing port of the customer server. ")

	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/healthz", serviceutils.Healthz)
	r.HandleFunc("/kill", serviceutils.Kill)
	r.HandleFunc("/case/list", listCases)
	r.HandleFunc("/case/register/{caseid}", registerCase)

	http.Handle("/", r)
	log.Println("listening on :", *http_port)
	go http.ListenAndServe(*http_port, nil)

	go clockclient.CreateClockClient(tock)

	lis, err := net.Listen("tcp", *rpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, newCustomerServer())
	log.Printf("Customer started at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// api_conn, _ := grpc.Dial("localhost:8002")
	// def api_conn.Close()

}
