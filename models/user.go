package models

import "net/url"

type User struct {
	Email   string
	SelfURL url.URL
}
