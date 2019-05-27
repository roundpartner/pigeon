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

	log.Printf("[INFO] [%s] Server starting on port %d", ServiceName, port)
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
	Check(rs.Router)
	rs.Router.HandleFunc("/email", rs.SendEmail).Methods("POST")
	rs.Router.HandleFunc("/template", rs.ViewTemplate).Methods("POST")
	rs.Router.HandleFunc("/verify", rs.VerifyAddress).Methods("POST")
	rs.Mail = NewMailService()
	return rs
}

func (rs *RestServer) SendEmail(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Received request to send email", ServiceName)
	decoder := json.NewDecoder(req.Body)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		log.Printf("[ERROR] [%s] Bad Request: %s", ServiceName, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if msg.To == "" {
		log.Printf("[ERROR] [%s] Bad Request: To address is required", ServiceName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs.Mail.QueueEmail(msg)
	w.WriteHeader(http.StatusNoContent)
}

func (rs *RestServer) ViewTemplate(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Received request to view template", ServiceName)
	decoder := json.NewDecoder(req.Body)
	msg := &Message{}
	err := decoder.Decode(msg)
	if err != nil {
		log.Printf("[ERROR] [%s] Bad Request: %s\n", ServiceName, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rs.Mail.AssembleTemplate(msg)
	if err != nil {
		log.Printf("[ERROR] [%s] Bad Request: %s\n", ServiceName, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[ERROR] [%s] Marshal Error: %s\n", ServiceName, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

type Lookup struct {
	Email   string `json:"email,omitempty"`
	Ip      string `json:"ip"`
	Blocked bool   `json:"blocked"`
}

func (rs *RestServer) VerifyAddress(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] [%s] Received request to verify email", ServiceName)
	decoder := json.NewDecoder(req.Body)
	lookup := &Lookup{}
	err := decoder.Decode(lookup)
	if err != nil {
		log.Printf("[ERROR] Bad Request: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] [%s] Looking up email=\"%s\" ip=\"%s\"", ServiceName, lookup.Email, lookup.Ip)
	lookup.verify()

	buf, err := json.Marshal(lookup)
	if err != nil {
		log.Printf("[ERROR] [%s] Marshal Error: %s\n", ServiceName, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (lookup *Lookup) verify() {
	if false == lookup.Blocked && "" != lookup.Email {
		lookup.Blocked = !ValidateEmail(lookup.Email)
	}

	if false == lookup.Blocked && "" != lookup.Ip {
		lookup.Blocked = CheckBlackList(lookup.Ip)
	}
}
