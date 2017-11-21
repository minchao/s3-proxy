package s3proxy

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
)

func BasicAuthHandler(next http.Handler, auth BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			unauthorized(w, r)
			return
		}

		if !auth.authenticate(username, password) {
			unauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+r.Host+`"`)
	http.Error(w, "", http.StatusUnauthorized)
}

type BasicAuth interface {
	authenticate(username, password string) bool
}

type SimpleBasicAuth struct {
	Username string
	Password string
}

func (ba *SimpleBasicAuth) authenticate(username, password string) bool {
	return username == ba.Username && password == ba.Password
}

type GithubOrgAuth struct {
	LoginOrg string
}

func (ga *GithubOrgAuth) authenticate(username, password string) bool {
	tp := github.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client := github.NewClient(tp.Client())

	orgs, _, err := client.Organizations.List(context.Background(), "", nil)
	if err != nil {
		return false
	}

	for _, org := range orgs {
		if org.Login != nil && ga.LoginOrg == *org.Login {
			return true
		}
	}
	return false
}
