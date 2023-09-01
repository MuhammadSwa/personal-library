package controllers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/muhammadswa/personal-library/pkg/repositories"
)

type Controllers struct {
	booksRespsitory *repositories.BooksRepository
	usersRepository *repositories.UsersRepository
	session         *scs.SessionManager
}

func New(booksRespo *repositories.BooksRepository, usersRepo *repositories.UsersRepository,
	session *scs.SessionManager) *Controllers {

	return &Controllers{
		usersRepository: usersRepo,
		booksRespsitory: booksRespo,
		session:         session,
	}

}
