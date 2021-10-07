package main

import (
    "fmt"
    "net/http"
	"log"
	"github.com/gorilla/mux"
)


//Global Variable!!!!
var assignedCases []supportCase 


func hello(w http.ResponseWriter, req *http.Request) {

    fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

type supportCase struct {
	CaseID string
	State string //State is the current case state 
}

func calcNextTick() {
	//Take the first assigned case, and figure out what to do with it. 
	currentCase := assignedCases[0]
	log.Println("Current case: ", currentCase)
	return 
}

func caseAssign(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// var caseId, ok = req.URL.Query()["caseid"]
	// if !ok || len(caseId[0]) < 1 {
    //     log.Println("Url Param 'caseid' is missing")
    //     return 
    // }
	caseid := vars["caseid"]
	c := supportCase{CaseID: caseid, State: "Assigned"}
	assignedCases = append(assignedCases, c)
	log.Println("case assigned: ", caseid)
	fmt.Fprintf(w, "case accepted %s, assigned cases: %d \n", caseid, len(assignedCases))

	//Trigger next work item recalc 
	return 
}

// work on the next assigned thing. 
func tick(w http.ResponseWriter, req *http.Request) {

}

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/hello", hello)
    r.HandleFunc("/headers", headers)
	r.HandleFunc("/assign/{caseid}", caseAssign)
	r.HandleFunc("/tick", tick)

	http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}