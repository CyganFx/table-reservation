package service

import (
	"github.com/CyganFx/table-reservation/pkg/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
)

type user struct {
	repo UserRepo
}

type UserRepo interface {
	Create(name, email, mobile, password string) error
	GetById(id int) (*domain.User, error)
	Update(user *domain.User) error
	Authenticate(email, password string) (int, error)
}

func NewUser(repo UserRepo) *user {
	return &user{repo: repo}
}

func (u *user) Save(form *forms.Form) (bool, error) {
	form.Required("name", "email", "mobile", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 5)
	form.MaxLength("name", 50)
	form.MaxLength("email", 100)

	if !form.Valid() {
		return false, nil
	}

	err := u.repo.Create(form.Get("name"),
		form.Get("email"),
		form.Get("mobile"),
		form.Get("password"))
	if err != nil {
		return true, err
	}

	return true, nil
}

func (u *user) SignIn(email, password string) (int, error) {
	return u.repo.Authenticate(email, password)
}

func (u *user) FindById(id int) (*domain.User, error) {
	return u.repo.GetById(id)
}

func (u *user) Update(user *domain.User) error {
	return u.repo.Update(user)
}
