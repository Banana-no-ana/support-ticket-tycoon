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

	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/grpc"

	"os/exec"

	"github.com/gorilla/mux"
)

var cases map[int32]*pb.Case
var workers map[int]*Worker
var nextCaseId int32 = 1 //TODO: Change to use closure instead
var nextWorkerId int = 100
var customerConn pb.CustomerClient

type Scenario struct {
	Difficulty pb.Difficulty
	Workers    []Worker `json:"workers"`
}

type Worker struct {
	WorkerID   int             //Assigned worker ID
	Name       string          //generateCased from a list of names
	FaceID     int             //Icon for worker face
	Skills     pb.WorkerSkill  //Worker's skills
	Connection pb.WorkerClient //Client connection to the worker.
}

type AssignRequest struct {
	WorkerID int
	CaseID   int
}

func generateCase(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to create a new case")
	c := pb.Case{CaseID: nextCaseId, Status: pb.CaseStatus_New, Assignee: 0,
		CustomerID: rand.Int31()%5 + 1, CustomerSentiment: 3}
	cases[nextCaseId] = &c
	log.Println("case created: ", nextCaseId)
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal CaseID: ", nextCaseId)
	}

	if customerConn != nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		log.Println("Registering case ", c.CaseID, " with the customer service")
		customerConn.RegisterCase(ctx, &c)
	}

	fmt.Fprintf(w, string(b))
	nextCaseId++
}

func getCaseUpdate(c *pb.Case) *pb.Case {
	if c.Assignee == 0 {
		//Case is not yet assigned.
		return c
	}
	w := workers[int(c.Assignee)]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	newc, err := w.Connection.GetCaseState(ctx, c)

	if err != nil {
		log.Println("Unable to get case update. Error: ", err.Error())
	}

	return newc
}

func listCases(w http.ResponseWriter, req *http.Request) {

	caseArray := []string{}

	for _, c := range cases {
		*c = *getCaseUpdate(c)
		d, _ := protojson.Marshal(c)
		caseArray = append(caseArray, string(d))
	}

	b, err := json.Marshal(caseArray)
	if err != nil {
		log.Fatal("Failed to marshal the list f cases")
	}

	fmt.Fprintf(w, string(b))

	// example := &pb.Case{CaseID: 1000}
	// jsonBytes, _ := protojson.Marshal(example)
	// fmt.Fprintf(w, string(jsonBytes))

}

func workerClientConnection(w *Worker) pb.WorkerClient {
	if w.Connection != nil {
		return w.Connection
	} else {
		createWorkerClient(w)
		return w.Connection
	}
}

func assign(w http.ResponseWriter, req *http.Request) {
	var ar AssignRequest
	err := json.NewDecoder(req.Body).Decode(&ar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := workers[ar.WorkerID]; !ok {
		fmt.Println("404 Failed: Worker ", ar.WorkerID, " is not found, not assigned")
		return
	}

	form_caseID32 := int32(ar.CaseID)

	log.Println("Received request to assign case ", ar.CaseID, "to worker ", ar.WorkerID)
	// wtfForm, _ := req.FormValue(("caseid"))
	// log.Println(wtfForm)

	c := cases[int32(ar.CaseID)]
	// *c = *getCaseUpdate(&pb.Case{CaseID: int32(ar.CaseID)})

	//Unassign from previous worker
	if c.Assignee != int32(0) {
		log.Println("Unassigning case: ", form_caseID32, " from worker: ", c.Assignee)
		curAssigneeConn := workerClientConnection(workers[int(c.Assignee)])
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		UnAssignResult, _ := curAssigneeConn.Unassign(ctx, &pb.Case{CaseID: form_caseID32})
		if !UnAssignResult.Success {
			log.Println("404: UnAssigning case ", form_caseID32, " from worker: ", ar.WorkerID, " failed")
			fmt.Fprintf(w, "404: UnAssigning case %d from worker %d failed", form_caseID32, ar.WorkerID)
			return
		}
	}

	//Assign to the worker
	clientConn := workerClientConnection(workers[ar.WorkerID])
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.Assignee = int32(ar.WorkerID)
	assignResult, err := clientConn.Assign(ctx, c)
	if assignResult.GetSuccess() {
		log.Println("Assigning case ", form_caseID32, " was successful")
		fmt.Fprintf(w, "Assigning case %d was successful", form_caseID32)
	} else {
		log.Println("Assigning case ", form_caseID32, " failed")
		log.Println(err)
		fmt.Fprintf(w, "Assigning case %d failed", form_caseID32)
	}
}

func getcase(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	caseid, _ := strconv.Atoi(vars["caseid"])

	log.Println("Received GET request for case: ", caseid)
	curCase := cases[int32(caseid)]

	*curCase = *getCaseUpdate(curCase)

	caseJson, err := protojson.Marshal(curCase)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "400: Failed to marshal case: %d", caseid)
		return
	}
	fmt.Fprintf(w, string(caseJson))
}

func registerWorker(w http.ResponseWriter, req *http.Request) {
	//TODO implement registering worker with the API server.
}

func listWorkers(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to list workers ")

	workerArray := []Worker{}

	for _, w := range workers {
		workerArray = append(workerArray, *w)
	}

	b, err := json.Marshal(workerArray)
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

func createOrReplaceWorker(w *Worker) {
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
	workers[w.WorkerID] = w
}

func createWorkerRequest(w http.ResponseWriter, req *http.Request) {

	//Create a new worker.
	log.Println("Starting worker: ", nextWorkerId)
	worker := Worker{WorkerID: nextWorkerId, FaceID: rand.Intn(80), Name: "Unnamed"}
	// n := "-worker_id=" + strconv.Itoa(nextWorkerId)

	createOrReplaceWorker(&worker)
	nextWorkerId++

	fmt.Fprintf(w, "Created worker \n")

	logmsg, _ := json.Marshal(worker)
	fmt.Fprintf(w, string(logmsg))
}

//Add a worker into the list for testing. Takes a POST request.
func addWorkerRequest(w http.ResponseWriter, req *http.Request) {
	var ww Worker
	err := json.NewDecoder(req.Body).Decode(&ww)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Adding worker: ", ww.WorkerID)
	workers[ww.WorkerID] = &ww
	// wtfForm, _ := req.FormValue(("caseid"))
	// log.Println(wtfForm)

	createWorkerClient(&ww)
	fmt.Fprintf(w, "Added worker %d", ww.WorkerID)
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

func createCustomerService() {
	if customerConn != nil {
		return
	}

	serviceaddr := "http://localhost:8005"
	_, err := http.Get(serviceaddr + "/healthz")
	if err != nil {
		log.Println("Creating the customer service ")
		cmd := exec.Command("go", "run", "customer.go")

		err = cmd.Start()
		// err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	connectToCustomerService()

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
		cur.FaceID = rand.Intn(80)
		createOrReplaceWorker(&cur)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go cur.Connection.SetWorkerSkills(ctx, &cur.Skills)
		numCreatedWorkers++
	}
	fmt.Fprintf(w, "Created "+strconv.Itoa(numCreatedWorkers)+" Workers \n")

	createCustomerService()

	fmt.Fprintf(w, "Loaded scenario from: "+f+"\n")

	log.Println("Finished loading scenario from: " + f + "\n")

}

func destroyWorker(w *Worker) {
	log.Println("Destroying worker: ", strconv.Itoa(w.WorkerID))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.Connection.KillWorker(ctx, &pb.Response{Success: true})
}

func unloadScenario(w http.ResponseWriter, req *http.Request) {
	log.Println("Received Request to unload scenario. Destroying all workers")

	for _, w := range workers {
		destroyWorker(w)
		delete(workers, w.WorkerID)
	}
	cases = make(map[int32]*pb.Case)

	fmt.Fprintf(w, "Removed all workers and cases")
}

func RegisterCustomerService(w http.ResponseWriter, req *http.Request) {
	log.Println("Manually adding customer serivce running on port:8005/6")

	connectToCustomerService()

	fmt.Println(w, "Customer service connected")
}

func main() {
	workers = make(map[int]*Worker)
	cases = make(map[int32]*pb.Case)
	createCustomerService()

	r := mux.NewRouter()
	r.HandleFunc("/healthz", serviceutils.Healthz) //Assign a case to a worker. Must be a post request
	r.HandleFunc("/kill", serviceutils.Kill)

	r.HandleFunc("/case/create", generateCase)   //generateCase new case and return a case ID.
	r.HandleFunc("/case/generate", generateCase) //generateCase new case and return a case ID.
	r.HandleFunc("/case/assign", assign)         //Assign a case to a worker. Must be a post request
	r.HandleFunc("/case/list", listCases)        //List all cases and their statuses
	r.HandleFunc("/case/get/{caseid}", getcase)  //get info of a case.

	// r.HandleFunc("/worker/register", registerWorker) //register a worker. Don't need this
	r.HandleFunc("/worker/list", listWorkers)           // Expected to be called by the frontend to list all the workers.
	r.HandleFunc("/worker/create", createWorkerRequest) // Create workers
	r.HandleFunc("/worker/add", addWorkerRequest)       // Add worker created elsewhere.

	r.HandleFunc("/customer/register", RegisterCustomerService) // Add worker created elsewhere.

	r.HandleFunc("/scenario/load/{scenarioid}", loadScenario) // Create workers
	r.HandleFunc("/scenario/unload", unloadScenario)          // destroy all the existing workers and clear out all the cases.

	http.Handle("/", r)
	log.Println("Listening on ", ":8001")
	http.ListenAndServe(":8001", nil) //HTTP endpoint: Mostly used for frontend

}
