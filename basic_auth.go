package s3proxy

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
