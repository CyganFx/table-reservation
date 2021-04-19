package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/CyganFx/table-reservation/pkg/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const UniqueViolationCode = "23505"

type user struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *user {
	return &user{db: db}
}

func (u *user) Create(name, email, mobile, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.New("failed to generate hashed password")
	}
	query := `INSERT INTO users (name, email, mobile, password, created)
	VALUES($1, $2, $3, $4, $5)`

	_, err = u.db.Exec(context.Background(), query, name, email, mobile, string(hashedPassword), time.Now())
	if err != nil {
		postgresError := err.(*pgconn.PgError)
		if errors.As(err, &postgresError) {
			if postgresError.Code == UniqueViolationCode &&
				strings.Contains(postgresError.Message, "users_uc_email") {
				return domain.ErrDuplicateEmail
			}
		}
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

func (u *user) GetById(id int) (*domain.User, error) {
	query := `SELECT name, email, mobile, created FROM users WHERE id = $1`
	user := &domain.User{}

	err := u.db.QueryRow(context.Background(), query, id).
		Scan(&user.Name, &user.Email,
			&user.Mobile, &user.Created)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, domain.ErrNoRecord
		} else {
			return nil, fmt.Errorf("failed to make select statement: %v", err)
		}
	}

	user.ID = id

	return user, nil
}

func (u *user) Update(user *domain.User) error {
	query := `UPDATE users SET name = $2, mobile = $3 FROM users WHERE id = $1`

	_, err := u.db.Exec(context.Background(), query, user.ID, user.Name, user.Mobile)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("failed to update: %v", err)
	}
	return nil
}

func (u *user) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := "SELECT id, password FROM users WHERE email = $1"
	row := u.db.QueryRow(context.Background(), query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return -1, domain.ErrInvalidCredentials
		}
		return -1, fmt.Errorf("failed to scan: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, domain.ErrInvalidCredentials
		} else {
			return -1, fmt.Errorf("failed to compare hash and password: %v", err)
		}
	}
	return id, nil
}
