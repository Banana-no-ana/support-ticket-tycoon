package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Banana-no-ana/support-ticket-tycoon/backend/protos"
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

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/register/{port}", register)

	http.Handle("/", r)
	log.Println("Starting master clock on port 8000")
	startClock()
	go http.ListenAndServe(":8000", nil)

	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
