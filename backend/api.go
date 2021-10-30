///API server is the primary interaction with frontend interms of "doing something"
//Frontend says: assign cases from A to B, that's done by the API server.
//API server keeps all the state that frontend may need
//API server also responsible for adjusting game speed (controlled by tick rate)

//Manually generateCase a case. Expose both RPC and REST endpoint, but they end up calling the same function.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"context"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
	serviceutils "github.com/Banana-no-ana/support-ticket-tycoon/backend/utils"
	"google.golang.org/grpc"

	"os/exec"

	"github.com/gorilla/mux"
)

var cases []Case
var workers []*Worker
var nextCaseId int //TODO: Change to use closure instead
var nextWorkerId int = 100

type Scenario struct {
	Workers []Worker `json:"workers"`
}

type Worker struct {
	WorkerID   int             //Assigned worker ID
	Name       string          //generateCased from a list of names
	FaceID     int             //Icon for worker face
	Skills     pb.WorkerSkill  //Worker's skills
	Connection pb.WorkerClient //Client connection to the worker.
}

type Case struct {
	CaseID            int
	State             string //State is the current case state
	Assignee          int    //worker UID.
	CustomerID        int    //Which customer is it
	CustomerSentiment int    //Customer's current sentiment of the case (range between 1, 2, 3, 4, 5 (5 being happy))
}

func generateCase(w http.ResponseWriter, req *http.Request) {
	c := Case{CaseID: nextCaseId, State: "New", Assignee: 0,
		CustomerID: rand.Intn(5) + 1, CustomerSentiment: 3}
	cases = append(cases, c)
	log.Println("case created: ", nextCaseId)
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal caseID: ", nextCaseId)
	}
	fmt.Fprintf(w, string(b))
	nextCaseId++
}

func listCases(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(cases)
	if err != nil {
		log.Fatal("Failed to marshal caseID: ", nextCaseId)
	}
	fmt.Fprintf(w, string(b))
}

func getWorkerfromID(id int) *Worker {
	for _, x := range workers {
		if x.WorkerID == id {
			return x
		}
	}
	return &Worker{}
}

func assign(w http.ResponseWriter, req *http.Request) {
	//TODO: Unassign from current assignee

	if err := req.ParseForm(); err != nil {
		log.Println(w, "ParseForm() err: %v", err)
		return
	}

	form_caseID, _ := strconv.Atoi(req.FormValue("caseid"))
	form_caseID32 := int32(form_caseID)
	form_workerID, _ := strconv.Atoi(req.FormValue("workerid"))

	log.Println("Received request to assign case ", form_caseID, "to worker ", form_workerID)
	// wtfForm, _ := req.FormValue(("caseid"))
	// log.Println(wtfForm)

	clientConn := getWorkerfromID(form_workerID).Connection

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result, err := clientConn.Assign(ctx, &pb.Case{CaseID: form_caseID32})
	if result.GetSuccess() {
		log.Println("Assigning case ", form_caseID32, " was successful")
		fmt.Fprintf(w, "Assigning case %d was successful", form_caseID32)
	} else {
		log.Println("Assigning case ", form_caseID32, " failed")
		log.Println(err)
		fmt.Fprintf(w, "Assigning case %d failed", form_caseID32)
	}
}

func getcase(w http.ResponseWriter, req *http.Request) {
	//Unassign from current assignee
	//Assign to new assignee

}

func registerWorker(w http.ResponseWriter, req *http.Request) {
	//TODO implement registering worker with the API server.
}

func listWorkers(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(workers)
	if err != nil {
		log.Fatal("Failed to marshal caseID: ", nextCaseId)
	}
	fmt.Fprintf(w, string(b))
}

func createWorkerClient(w *Worker) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	workeraddr := "localhost:" + strconv.Itoa(10000+w.WorkerID)
	worker_conn, _ := grpc.Dial(workeraddr, opts...)
	// defer worker_conn.Close()

	w.Connection = pb.NewWorkerClient(worker_conn)
}

func createWorker(w *Worker) {
	rpc_port := "-rpc_port=:" + strconv.Itoa(10000+w.WorkerID)
	http_port := "-http_port=:" + strconv.Itoa(9000+w.WorkerID)

	log.Println("Creating worker: ", w.WorkerID, rpc_port, http_port)
	cmd := exec.Command("go", "run", "worker.go", "-worker_id="+strconv.Itoa(w.WorkerID), rpc_port, http_port)

	err := cmd.Start()
	// err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	createWorkerClient(w)
	workers = append(workers, w)

}

func createWorkerRequest(w http.ResponseWriter, req *http.Request) {

	//Create a new worker.
	log.Println("Starting worker: ", nextWorkerId)
	worker := Worker{WorkerID: nextWorkerId, FaceID: 1, Name: "Unnamed"}
	// n := "-worker_id=" + strconv.Itoa(nextWorkerId)

	createWorker(&worker)
	nextWorkerId++

	fmt.Fprintf(w, "Created worker \n")

	logmsg, _ := json.Marshal(worker)
	fmt.Fprintf(w, string(logmsg))
}

func loadScenario(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	scenarioid := vars["scenarioid"]

	//TODO: Need to change to some sort of absolute pathing instead of relative pathing
	f := "../scenarios/s" + scenarioid + ".json"

	log.Println("Loading scenario from: ", f)
	sceanriofile, err := os.Open(f)
	if err != nil {
		fmt.Println(err)
	}
	defer sceanriofile.Close()

	bytevalue, _ := ioutil.ReadAll(sceanriofile)
	var scenario Scenario

	err = json.Unmarshal(bytevalue, &scenario)
	if err != nil {
		fmt.Println("error:", err)
	}

	var numCreatedWorkers int = 0
	for _, worker := range scenario.Workers {
		//TODO: What do we do if these workers exist? Kill their existing work probably.
		cur := Worker(worker)
		createWorker(&cur)
		numCreatedWorkers++
	}

	fmt.Fprintf(w, "Loaded scenario from: "+f+"\n")
	fmt.Fprintf(w, "Created "+strconv.Itoa(numCreatedWorkers)+" Workers \n")

}

func destroyWorker(w *Worker) {
	log.Println("Destoring worker: ", strconv.Itoa(w.WorkerID))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Connection.KillWorker(ctx, &pb.Response{Success: true})
}

func unloadScenario(w http.ResponseWriter, req *http.Request) {
	log.Println("Received Request to unload scenario. Destroying all workers")

	for _, w := range workers {
		destroyWorker(w)
	}
	workers = nil

	fmt.Fprintf(w, "Removed all workers")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", serviceutils.Healthz) //Assign a case to a worker. Must be a post request
	r.HandleFunc("/kill", serviceutils.Kill)

	r.HandleFunc("/case/create", generateCase)  //generateCase new case and return a case ID.
	r.HandleFunc("/case/assign", assign)        //Assign a case to a worker. Must be a post request
	r.HandleFunc("/case/list", listCases)       //List all cases and their statuses
	r.HandleFunc("/case/get/{caseid}", getcase) //get info of a case.

	// r.HandleFunc("/worker/register", registerWorker) //register a worker. Don't need this
	r.HandleFunc("/worker/list", listWorkers)           // Expected to be called by the frontend to list all the workers.
	r.HandleFunc("/worker/create", createWorkerRequest) // Create workers

	r.HandleFunc("/scenario/load/{scenarioid}", loadScenario) // Create workers
	r.HandleFunc("/scenario/unload", unloadScenario)          // destroy all the existing workers and clear out all the cases.

	http.Handle("/", r)
	log.Println("Listening on ", ":8001")
	http.ListenAndServe(":8001", nil) //HTTP endpoint: Mostly used for frontend

}
