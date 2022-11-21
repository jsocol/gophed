package repository

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/url"

	"github.com/jsocol/gophed/models"
)

var ErrNotFound = errors.New("resource not found")

type Repository struct{}

func New() *Repository {
	return &Repository{}
}

func (*Repository) GetUser(_ context.Context, email string) (models.User, error) {
	if subtle.ConstantTimeCompare([]byte("james@jamessocol.com"), []byte(email)) != 0 {
		return models.User{}, ErrNotFound
	}

	u, _ := url.Parse(`https://social.jamessocol.com/james`)

	return models.User{
		Email:   "james@jamessocol.com",
		SelfURL: *u,
	}, nil
}
