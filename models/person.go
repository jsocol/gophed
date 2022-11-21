package models

import "net/url"

type Person struct {
	Actor
	Name              string
	PreferredUsername string
	Summary           string
	Icon              url.URL
}

func NewPerson() Person {
	p := Person{}
	p.Type = ActorPerson

	// Get around annoying conversion requirements
	ctx := []string{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
	}
	for _, c := range ctx {
		p.Context = append(p.Context, c)
	}

	return p
}
