package main

import "net/http"

func (rs *RestServer) Check(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
