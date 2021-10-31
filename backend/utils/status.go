package serviceutils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Healthz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok")
}

func Kill(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request to terminate")
	go Kill2()
}

func Kill2() {
	log.Println("Received request to terminate")
	time.Sleep(50 * time.Millisecond)
	os.Exit(0)
}
