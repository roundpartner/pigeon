package main

import (
	"log"
	"net"
	"strings"
)

func ValidateEmail(email string) bool {
	components := strings.Split(email, "@")
	if len(components) != 2 {
		log.Printf("[ERROR] [%s] %s is not a valid email address", ServiceName, email)
		return false
	}
	mxrecords, err := net.LookupMX(components[1])
	if err != nil {
		log.Printf("[ERROR] [%s] %s", ServiceName, err.Error())
		return false
	}
	return len(mxrecords) > 0
}
