package main

import (
	"github.com/mrichman/godnsbl"
)

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
