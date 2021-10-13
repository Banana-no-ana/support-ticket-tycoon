package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/protos/clock"
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

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/register/{port}", register)

	http.Handle("/", r)
	log.Println("Starting master clock on port 8000")
	startClock()
	http.ListenAndServe(":8000", nil)
}
