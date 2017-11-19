package s3proxy

type BasicAuth interface {
	authenticate(username, password string) bool
}

type SimpleBasicAuth struct {
	Username string
	Password string
}

func (a *SimpleBasicAuth) authenticate(username, password string) bool {
	if username == a.Username && password == a.Password {
		return true
	}
	return false
}
