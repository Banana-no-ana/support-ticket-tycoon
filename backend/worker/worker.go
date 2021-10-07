package main

import (
    "fmt"
    "net/http"
)

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
	State string
}

func caseAssign(w http.ResponseWriter, req *http.Request, c []) {


}

func caseAssign(w http.ResponseWriter, req *http.Request, ) {


}



func main() {

	cases := make([]supportCase, 20)

	
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)

	http.handleFunc("/assign", caseAssign, cases)
	http.handleFunc("/caseList", caseList, cases)	

    http.ListenAndServe(":8080", nil)
}