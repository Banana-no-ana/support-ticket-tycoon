//Customers respond to case messages and may assign csat to close case.
//Customers do not create cases. That's done by the case generator
//API polls the customer on their sentiment.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Banana-no-ana/support-ticket-tycoon/backend/clockclient"
	"github.com/gorilla/mux"
)

//When the clock ticks, it tocks us.
func tock() {
	log.Println("We've been Tock'd ")
}

func caseStatus(w http.ResponseWriter, req *http.Request) {

}

func main() {
	http_port := flag.String("http_port", ":8020", "set the http listneing port of the customer server. ")
	// rpc_port := flag.String("rpc_port", ":8021", "set the rpc listneing port of the customer server. ")

	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/case/{caseid}", caseStatus)

	http.Handle("/", r)
	log.Println("listening on :", *http_port)
	go http.ListenAndServe(*http_port, nil)

	clockclient.CreateClockClient(tock)

	// lis, err := net.Listen("tcp", *rpc_port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterWorkerServer(s, newWorkerServer())
	// log.Printf("Worker %d started at %v", workerID, lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	// api_conn, _ := grpc.Dial("localhost:8002")
	// def api_conn.Close()

}
