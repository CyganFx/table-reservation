package service

import (
	"errors"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/domain"
	"github.com/CyganFx/table-reservation/ez-booking/pkg/validator/forms"
	"golang.org/x/crypto/bcrypt"
)

const UserRoleId = 2

type user struct {
	repo UserRepo
}

type UserRepo interface {
	Create(name, email, mobile, hashedPassword string, roleId int) error
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
	form.MinLength("mobile", 11)
	form.MaxLength("mobile", 12)
	form.MaxLength("name", 50)
	form.MaxLength("email", 100)

	if !form.Valid() {
		return false, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Get("password")), 12)
	if err != nil {
		return false, errors.New("failed to generate hashed password")
	}

	err = u.repo.Create(form.Get("name"),
		form.Get("email"),
		form.Get("mobile"),
		string(hashedPassword),
		UserRoleId,
	)
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
