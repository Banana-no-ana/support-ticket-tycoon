package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func registerwithClock() {
	//TODO: Figure out how to actually do this.
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/assign/{caseid}", caseAssign)
	r.HandleFunc("/tick/{ticknum}", tick)

	http.Handle("/", r)

	registerwithClock()
	http.ListenAndServe(":8080", nil)
}
