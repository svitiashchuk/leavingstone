package sqlite

import (
	"database/sql"
	"leavingstone"

	_ "github.com/mattn/go-sqlite3"
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

func (us *UserService) FindByID(id int) (*leavingstone.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token, password FROM users WHERE id = ?", id)

	user := &leavingstone.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Find(email string) (*leavingstone.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token, password FROM users WHERE email = ?", email)

	user := &leavingstone.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) FindByToken(token string) (*leavingstone.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token FROM users WHERE token = ?", token)

	user := &leavingstone.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) AllUsers() ([]*leavingstone.User, error) {
	uu := []*leavingstone.User{}

	rows, err := us.db.Query(`
		SELECT u.id, u.name, u.email, u.token, u.start, u.extra_vacation, l.id, l.start, l.end, l.type, l.approved, l.user_id
		FROM users u
		INNER JOIN leaves l ON u.id = l.user_id
	`)

	if err != nil {
		return nil, err
	}

	// Map to store users by ID and collect leaves as child field for each user
	users := map[int]*leavingstone.User{}

	for rows.Next() {
		user := leavingstone.User{}
		leave := leavingstone.Leave{}

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
			&leave.UserID,
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
