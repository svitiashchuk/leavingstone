package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DSN = "file:database.db?cache=shared&mode=rwc"

type User struct {
	Name  string
	Email string
	Token string
}

type UserService struct {
	db *sql.DB
}

func NewUserService() (*UserService, error) {
	db, err := sql.Open("sqlite3", DSN)
	if err != nil {
		return nil, err
	}

	return &UserService{db}, nil
}

func (us *UserService) User(email string) (*User, error) {
	row := us.db.QueryRow("SELECT name, email, token FROM users WHERE email = ?", email)

	user := &User{}
	err := row.Scan(&user.Name, &user.Email, &user.Token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) AllUsers() ([]*User, error) {
	uu := []*User{}

	rows, err := us.db.Query("SELECT name, email, token FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := &User{}

		err := rows.Scan(&user.Name, &user.Email, &user.Token)
		if err != nil {
			return nil, err
		}

		uu = append(uu, user)
	}

	return uu, nil
}
