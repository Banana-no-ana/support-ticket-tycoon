package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var services []int
var tick int = 0

func notifySubscribers() {
	for _, s := range services {
		//do nothing for now
		serv := fmt.Sprintf("http://localhost:%d/tick/%d", s, tick)
		log.Println("notifying: ", serv)
		go http.Get(serv)
	}
}

func startClock() {

	ticker := time.NewTicker(2000 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				tick = tick + 1
				log.Println("Ticking: ", tick)
				notifySubscribers()
			}
		}
	}()

}

func register(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	port, _ := strconv.Atoi(vars["port"])

	log.Println("Received request to register on port: ", port)
	fmt.Fprintf(w, "hello, %d \n", port)
	services = append(services, port)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/register/{port}", register)

	http.Handle("/", r)
	log.Println("Starting master clock on port 8000")
	startClock()
	http.ListenAndServe(":8000", nil)
}
