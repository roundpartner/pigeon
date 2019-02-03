package main

import (
	"net/http"
	"testing"
)

func TestShutdownGracefully(t *testing.T) {
	server := &http.Server{}
	ShutdownGracefully(server)
}
