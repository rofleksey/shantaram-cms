package util

type ContextKey string

func (c ContextKey) String() string {
	return "shantaram_" + string(c)
}

var UsernameContextKey ContextKey = "username"
var IpContextKey ContextKey = "ip"
