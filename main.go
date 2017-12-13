package main

import (
	"log"
	"os"
	"github.com/artyom/autoflags"
	"flag"
)

var ServerConfig = struct {
	Port int `flag:"port,port number to listen on"`
}{Port: 3411}


func main() {
	autoflags.Define(&ServerConfig)
	flag.Parse()

	log.SetOutput(os.Stdout)

	ListenAndServe(ServerConfig.Port)
}
