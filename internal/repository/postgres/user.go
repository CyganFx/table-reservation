package postgres

import (
	"context"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
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

func (u *user) Create(name, email, mobile, hashedPassword string, roleId int) error {
	query := `INSERT INTO users (name, role_id, email, mobile, password, created)
	VALUES($1, $2, $3, $4, $5, $6)`

	_, err := u.db.Exec(context.Background(), query, name, roleId, email, mobile, hashedPassword, time.Now())
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
	query := `SELECT u.name, u.role_id, r.name, u.email, u.mobile, u.created, u.profile_image_url
			FROM users u join roles r on u.role_id = r.id WHERE u.id = $1`
	user := domain.NewUser()

	err := u.db.QueryRow(context.Background(), query, id).
		Scan(&user.Name, &user.Role.ID, &user.Role.Name, &user.Email,
			&user.Mobile, &user.Created, &user.ImageURL)
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

func (u *user) SetProfileImage(filePath string, userID int) error {
	query := `UPDATE users SET profile_image_url = $1 WHERE id = $2`

	_, err := u.db.Exec(context.Background(), query, filePath, userID)
	if err != nil {
		return fmt.Errorf("failed to update profile image %v", err)
	}

	return nil
}

func (u *user) UpdateUserRoleByID(userID, roleID int) error {
	query := `UPDATE users SET role_id = $2 WHERE id = $1`

	_, err := u.db.Exec(context.Background(), query, userID, roleID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("failed to update: %v", err)
	}

	return nil
}

func (u *user) Query(cafeID int) ([]domain.User, error) {
	query := `SELECT id, name, email, mobile FROM users where 
			role_id != 1 and role_id != 3
			and 
			id not in
			(select user_id from blacklist where cafe_id = $1)`

	rows, err := u.db.Query(context.Background(), query, cafeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uu []domain.User

	for rows.Next() {
		u := usersPool.Get().(*domain.User)
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Mobile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign values to location struct from row")
		}
		uu = append(uu, *u)

		*u = domain.User{}
		usersPool.Put(u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return uu, nil
}

func (u *user) Report(userID, cafeID int) error {
	query := `INSERT INTO blacklist VALUES($1, $2)`

	_, err := u.db.Exec(context.Background(), query, userID, cafeID)
	if err != nil {
		return errors.Wrap(err, "inserting user in the blacklist")
	}

	return nil
}
