package main

import (
	"github.com/adtac/go-akismet/akismet"
	"github.com/mrichman/godnsbl"
	"log"
	"os"
)

func CheckBlackList(ip string) bool {
	lookup := Lookup{Ip: ip}
	blacklists := []string{
		"sbl.spamhaus.org",
		"xbl.spamhaus.org",
		"spam.spamrats.com",
		"all.s5h.net",
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

func CheckAkismet(msg *Message) bool {
	akismetKey, isSet := os.LookupEnv("AKISMET_KEY")
	if isSet == false {
		log.Printf("[INFO] AKISMET_KEY is not set")
		return false
	}
	isSpam, err := akismet.Check(&akismet.Comment{
		Blog:               msg.Website,
		UserIP:             msg.SenderIp,
		UserAgent:          msg.UserAgent,
		CommentType:        "contact-form",
		CommentAuthor:      msg.FromName,
		CommentAuthorEmail: msg.From,
		CommentContent:     msg.Text,
	}, akismetKey)

	if err != nil {
		log.Printf("[ERROR] Akismet error: %s", err.Error())
		return false
	}
	return isSpam
}
