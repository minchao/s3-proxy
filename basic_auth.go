package s3proxy

import (
	"context"

	"github.com/google/go-github/github"
)

type BasicAuth interface {
	authenticate(username, password string) bool
}

type SimpleBasicAuth struct {
	Username string
	Password string
}

func (ba *SimpleBasicAuth) authenticate(username, password string) bool {
	if username == ba.Username && password == ba.Password {
		return true
	}
	return false
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
