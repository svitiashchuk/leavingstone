package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"ptocker"
)

const DSN = "file:database.db?cache=shared&mode=rwc"

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

func (us *UserService) Find(email string) (*ptocker.User, error) {
	row := us.db.QueryRow("SELECT name, email, token FROM users WHERE email = ?", email)

	user := &ptocker.User{}
	err := row.Scan(&user.Name, &user.Email, &user.Token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) AllUsers() ([]*ptocker.User, error) {
	uu := []*ptocker.User{}

	rows, err := us.db.Query(`
		SELECT u.id, u.name, u.email, u.token, u.start, u.extra_vacation, l.id, l.start, l.end, l.type, l.approved
		FROM users u
		INNER JOIN leaves l ON u.id = l.user_id
	`)

	if err != nil {
		return nil, err
	}

	// Map to store users by ID and collect leaves as child field for each user
	users := map[int]*ptocker.User{}

	for rows.Next() {
		user := ptocker.User{}
		leave := ptocker.Leave{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Token,
			&user.Started,
			&user.ExtraVacation,
			&leave.ID,
			&leave.Start,
			&leave.End,
			&leave.Type,
			&leave.Approved,
		)

		if err != nil {
			return nil, err
		}

		if _, ok := users[user.ID]; !ok {
			users[user.ID] = &user
		}

		users[user.ID].Leaves = append(users[user.ID].Leaves, &leave)
	}

	for _, user := range users {
		uu = append(uu, user)
	}

	return uu, nil
}
