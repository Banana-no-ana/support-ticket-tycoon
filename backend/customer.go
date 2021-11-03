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

var difficulty pb.Difficulty
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
	difficulty = *dif
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

func generateCaseStage(stage int32) *pb.CaseStage {
	return &pb.CaseStage{}
}

func (CustomerServer) CustomerReply(ctx context.Context, c *pb.Case) (*pb.Case, error) {
	//Check if they're waiting for our reply.
	dif := rand.Int31n(difficulty.MaxDifficulty-difficulty.MinDifficulty) + difficulty.MinDifficulty
	if c.Status == "Waiting for Customer Reply" {
		if c.CurrentStage < dif {
			c.CurrentStage = c.CurrentStage + 1
			curStage := generateCaseStage(c.CurrentStage)
			c.CaseStages = append(c.CaseStages, curStage)
			return c, nil
		} else {
			c.Status = "Closed"
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
