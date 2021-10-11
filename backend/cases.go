//Master list of cases and case status.

//Manually generate a case. Expose both RPC and REST endpoint, but they end up calling the same function.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var cases []Case
var nextCaseId int //TODO: Change to use closure instead

type Case struct {
	CaseID string
	State  string //State is the current case state
}

func generate(w http.ResponseWriter, req *http.Request) {
	c := Case{CaseID: caseid, State: "Assigned"}
	assignedCases = append(assignedCases, c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %s, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate", generate) //generate new case and return a case ID.

	http.Handle("/", r)

	registerwithClock()
	http.ListenAndServe(":8082", nil)
}
