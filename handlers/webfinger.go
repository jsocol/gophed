package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jsocol/gophed/models"
	"github.com/jsocol/gophed/repository"
)

const jsonType = "application/json"

type link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

type response struct {
	Subject string `json:"subject"`
	Links   []link `json:"links"`
}

func newResponse() response {
	return response{
		Links: []link{
			{
				Rel:  "self",
				Type: "application/activity+json",
			},
		},
	}
}

type userGetter interface {
	GetUser(context.Context, string) (models.User, error)
}

type WebfingerHandler struct {
	Users userGetter
}

func (wh WebfingerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Add("Content-type", jsonType)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "bad http method"}`))
		return
	}

	params := r.URL.Query()
	resource := params.Get("resource")
	if resource == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "resource not specified"}`))
		return
	}
	if !strings.HasPrefix(resource, "acct:") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad resource name"}`))
		return
	}

	email := strings.TrimPrefix(resource, "acct:")
	subject, err := wh.Users.GetUser(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "resource not found"}`))
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	resp := newResponse()
	resp.Subject = resource
	resp.Links[0].Href = subject.SelfURL.String()

	respData, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}
