package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ListenAndServe(port int) {
	address := fmt.Sprintf(":%d", port)

	rs := NewRestServer()
	server := &http.Server{Addr: address, Handler: rs.Router}

	ShutdownGracefully(server)

	log.Printf("Server starting on port %d\n", port)
	err := server.ListenAndServe()
	if nil != err {
		log.Println(err.Error())
	}
}

type RestServer struct {
	Router            *mux.Router
	Mail              *MailService
	activeConnections int
}

func NewRestServer() *RestServer {
	rs := &RestServer{}
	rs.Router = mux.NewRouter()
	rs.Router.HandleFunc("/email", rs.SendEmail).Methods("POST")
	rs.Router.HandleFunc("/template", rs.ViewTemplate).Methods("POST")
	rs.Mail = NewMailService()
	return rs
}

func (rs *RestServer) SendEmail(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		log.Printf("[ERROR] Bad Request: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if msg.To == "" {
		log.Printf("[ERROR] Bad Request: To address is required\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs.Mail.QueueEmail(msg)
	w.WriteHeader(http.StatusNoContent)
}

func (rs *RestServer) ViewTemplate(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		log.Printf("[ERROR] Bad Request: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rs.Mail.AssembleTemplate(msg)
	if err != nil {
		log.Printf("[ERROR] Bad Request: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[ERROR] Marshal Error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
