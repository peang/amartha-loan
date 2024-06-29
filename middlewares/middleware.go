package middleware

import (
	"github.com/casbin/casbin/v2"
)

type Middleware struct {
	enforcer *casbin.Enforcer
}

func NewMiddleware(enfocer *casbin.Enforcer) *Middleware {
	return &Middleware{
		enforcer: enfocer,
	}
}
