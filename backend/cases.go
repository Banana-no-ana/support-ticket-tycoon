//Master list of cases and case status.

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate", generate) //generate new case and return a case ID.

	http.Handle("/", r)
	http.ListenAndServe(":8002", nil)
}
