package main

import (
	"fmt"
	"github.com/mrichman/godnsbl"
	"github.com/thomaslorentsen/go-spamc"
	"log"
	"time"
)

func CheckSpamAssassin(msg *Message) {
	t := time.Now()
	msgId := fmt.Sprintf("<%d.%s>", t.Unix(), msg.From)
	html := fmt.Sprintf("To: %s\n\rFrom: %s\n\rSubject: %s\n\rDate: %s\n\rMessage-ID: %s\n\r\n\r%s\n\r", msg.To, msg.From, msg.Subject, t.Format("Fri, 02 Jan 2006 15:04:05 -0700"), msgId, msg.Text)

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

func CheckBlackList(ip string) bool {
	lookup := Lookup{Ip: ip}
	blacklists := []string{
		"sbl.spamhaus.org",
		"xbl.spamhaus.org",
		"spam.spamrats.com",
	}
	for _, source := range blacklists {
		result := godnsbl.Lookup(source, lookup.Ip)
		if len(result.Results) > 0 {
			lookup.Blocked = result.Results[0].Listed
			if lookup.Blocked {
				return true
			}
		}
	}
	return false
}
