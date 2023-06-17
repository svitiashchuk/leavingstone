package sqlite

import (
	"database/sql"
	"fmt"

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

func (us *UserService) Get(email string) (*User, error) {
	row := us.db.QueryRow("SELECT name, email, token FROM users WHERE email = ?", email)

	user := &User{}
	err := row.Scan(&user.Name, &user.Email, &user.Token)
	if err != nil {
		return nil, err
	}

	fmt.Println("User: ", user)

	return user, nil
}
