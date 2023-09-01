package controllers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/muhammadswa/personal-library/pkg/repositories"
)

type Controllers struct {
	repos   *repositories.Repositories
	session *scs.SessionManager
}

func New(repos *repositories.Repositories,
	session *scs.SessionManager) *Controllers {

	return &Controllers{
		repos:   repos,
		session: session,
	}

}
