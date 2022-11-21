package models

import "net/url"

type ActorType interface {
	isActorType()
}

type actorType string

func (actorType) isActorType() {}

const ActorPerson actorType = "Person"

type Actor struct {
	ID        string
	Type      ActorType
	Context   []interface{}
	Inbox     url.URL
	Outbox    url.URL
	Following url.URL
	Followers url.URL
	Liked     url.URL
}
