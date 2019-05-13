package main

import (
	"flag"
	"github.com/artyom/autoflags"
	"log"
	"os"
)

var ServiceName = "pigeon"

var ServerConfig = struct {
	Port int `flag:"port,port number to listen on"`
}{
	Port: 3411,
}

func main() {
	log.SetOutput(os.Stdout)
	autoflags.Define(&ServerConfig)
	flag.Parse()

	serviceName, isSet := os.LookupEnv("SERVICE_NAME")
	if isSet {
		ServiceName = serviceName
	}

	ListenAndServe(ServerConfig.Port)
}
