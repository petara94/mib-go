package services

import "github.com/petara94/mib-go/internal/db"

type UserRepo interface {
	Create(user db.User) error
	GetByLogin(login string) (*db.User, error)
	Update(login string, user db.User) (*db.User, error)
	Delete(login string) error
}
