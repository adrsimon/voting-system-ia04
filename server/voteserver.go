package server

import (
	"fmt"
	"sync"
	"net/http"
	"time"
	"log"
)

type VoteServer struct {
	sync.Mutex
	id					string
	addr				string
}

func NewVoteServer(addr string) *VoteServer {
	return &VoteServer{id: addr, addr: addr}
}

func (vs *VoteServer) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}


func (vs *VoteServer) greet(w http.ResponseWriter, r *http.Request) {
	if !vs.checkMethod("GET", w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
	vs.Lock()
	defer vs.Unlock()

	msg := fmt.Sprintf("Hello world !")
	w.Write([]byte(msg))
}

func (vs *VoteServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", vs.greet)

	s := &http.Server{
		Addr:           vs.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening on", vs.addr)
	go log.Fatal(s.ListenAndServe())
}