package s3proxy

type BasicAuth interface {
	authenticate(username string, password string) bool
}

type SimpleBasicAuth struct {
	Username string
	Password string
}

func (a SimpleBasicAuth) authenticate(username string, password string) bool {
	if username == a.Username && password == a.Password {
		return true
	}
	return false
}
