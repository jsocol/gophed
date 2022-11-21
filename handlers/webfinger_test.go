package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsocol/gophed/handlers"
	"github.com/jsocol/gophed/models"
)

type mockGetter struct {
	user models.User
	err  error
}

func (m mockGetter) GetUser(context.Context, string) (models.User, error) {
	return m.user, m.err
}

func TestWebfinger_Success(t *testing.T) {
	u := mockGetter{
		user: models.User{
			Email: "test@example.com",
			SelfURL: url.URL{
				Scheme: "https",
				Host:   "localhost",
				Path:   "path",
			},
		},
	}

	wh := handlers.WebfingerHandler{
		Users: u,
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/.webfinger?resource=acct:test@example.com", nil)

	wh.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}
