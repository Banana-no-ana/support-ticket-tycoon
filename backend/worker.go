package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	"github.com/gorilla/mux"
)

//Global Variable!!!!
var assignedCases []Case

type Case struct {
	CaseID int
	State  string //State is the current case state
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

func caseAssign(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])
	c := Case{CaseID: caseid, State: "Assigned"}
	assignedCases = append(assignedCases, c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %d, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc
	return
}

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
	log.Printf("Connecting to the clock server")
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

func main() {
	listeningport := flag.String("listen_addr", ":8081", "set the listneing port of the worker. ")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/assign/{caseid}", caseAssign)
	r.HandleFunc("/unassign/{caseid}", caseAssign)
	r.HandleFunc("/tick/{ticknum}", tick)

	http.Handle("/", r)
	log.Println("listening on :", *listeningport)
	go http.ListenAndServe(*listeningport, nil)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, _ := grpc.Dial("localhost:8001", opts...)
	defer conn.Close()
	client := pb.NewClockClient(conn)
	registerwithClock(client)

}
