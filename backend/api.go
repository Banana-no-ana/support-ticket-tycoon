///API server is the primary interaction with frontend interms of "doing something"
//Frontend says: assign cases from A to B, that's done by the API server.
//API server keeps all the state that frontend may need
//API server also responsible for adjusting game speed (controlled by tick rate)

//Manually generate a case. Expose both RPC and REST endpoint, but they end up calling the same function.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var cases []Case
var nextCaseId int //TODO: Change to use closure instead

type Case struct {
	CaseID   int
	State    string //State is the current case state
	Assignee string //worker UID.
}

func generate(w http.ResponseWriter, req *http.Request) {
	c := Case{CaseID: nextCaseId, State: "New", Assignee: "Unassigned"}
	cases = append(cases, c)
	log.Println("case created: ", nextCaseId)
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Failed to marshal caseID: ", nextCaseId)
	}
	fmt.Fprintf(w, string(b))
	nextCaseId++
}

func list(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(cases)
	if err != nil {
		log.Fatal("Failed to marshal caseID: ", nextCaseId)
	}
	fmt.Fprintf(w, string(b))
}

func assign(w http.ResponseWriter, req *http.Request) {
	//Unassign from current assignee
	//Assign to new assignee
}

func getcase(w http.ResponseWriter, req *http.Request) {
	//Unassign from current assignee
	//Assign to new assignee
}

func registerWorker(w http.ResponseWriter, req *http.Request) {
	//TODO implement registering worker with the API server.
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/case/assign", assign)        //Assign a case to a worker. Must be a post request
	r.HandleFunc("/case/list", list)            //List all cases and their statuses
	r.HandleFunc("/case/create", generate)      //generate new case and return a case ID.
	r.HandleFunc("/case/get/{caseid}", getcase) //get info of a case.

	r.HandleFunc("/worker/register", registerWorker) //register a worker

	http.Handle("/", r)
	http.ListenAndServe(":8001", nil) //HTTP endpoint: Mostly used for frontend

}
