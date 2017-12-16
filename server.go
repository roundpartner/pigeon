package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ListenAndServe(port int) {
	address := fmt.Sprintf(":%d", port)

	rs := NewRestServer()
	server := &http.Server{Addr: address, Handler: rs.Router}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		signal.Stop(c)
		log.Println("Server shutting down gracefully")
		server.Shutdown(nil)
	}()

	log.Printf("Server starting on port %d\n", port)
	err := server.ListenAndServe()
	if nil != err {
		log.Println(err.Error())
	}
}

type RestServer struct {
	Router *mux.Router
	Mail   *MailService
}

func NewRestServer() *RestServer {
	rs := &RestServer{}
	rs.Router = mux.NewRouter()
	rs.Router.HandleFunc("/email", rs.SendEmail).Methods("POST")
	rs.Mail = NewMailService()
	return rs
}

func (rs *RestServer) SendEmail(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if msg.To == "" {
		log.Printf("Error: To address is required\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs.Mail.QueueEmail(msg)
	w.WriteHeader(http.StatusNoContent)
}
