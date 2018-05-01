package main

import (
	"fmt"
	"github.com/saintienn/go-spamc"
	"log"
)

func CheckSpamAssassin(msg *Message) {
	html := fmt.Sprintf("To: %s\n\rFrom: %s\n\rSubject: %s\n\rMessage: %s", msg.To, msg.From, msg.Subject, msg.Html)

	client := spamc.New("127.0.0.1:783", 10)
	reply, err := client.Report(html)
	if err != nil {
		msg.SpamScore = 1.0
		log.Println(reply, err)
		return
	}
	log.Println(reply.Vars["report"])
	msg.SpamScore = reply.Vars["spamScore"].(float64)
	msg.IsSpam = reply.Vars["isSpam"].(bool)
}
