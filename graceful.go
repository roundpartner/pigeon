package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var activeConnections = 0

func ShutdownGracefully(server *http.Server) {
	server.ConnState = func(conn net.Conn, state http.ConnState) {
		if "new" == state.String() {
			activeConnections++
		}
		if "closed" == state.String() {
			activeConnections--
		}
		if "hijacked" == state.String() {
			activeConnections--
		}
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		signal.Stop(c)

		serviceAvailable = false

		log.Printf("[INFO] [%s] Waiting for active connections to stop", ServiceName)
		for activeConnections > 0 {
			time.Sleep(time.Millisecond)
		}
		log.Printf("[INFO] [%s] Server shutting down gracefully", ServiceName)

		err := server.Shutdown(nil)
		if nil != err {
			log.Printf("[ERROR] [%s] Error shutting down server: %s", ServiceName, err.Error())
		}
	}()
}
