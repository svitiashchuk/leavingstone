package sqlite

import (
	"database/sql"
	"fmt"
	"leavingstone/internal/model"
	"strings"
	"time"

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

func (us *UserService) FindByID(id int) (*model.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token, password, start FROM users WHERE id = ?", id)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Password, &user.Started)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Find(email string) (*model.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token, password, start FROM users WHERE email = ?", email)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Password, &user.Started)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) FindByToken(token string) (*model.User, error) {
	row := us.db.QueryRow("SELECT id, name, email, token, start FROM users WHERE token = ?", token)

	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Started)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) AllUsers() ([]*model.User, error) {
	uu := []*model.User{}

	rows, err := us.db.Query(`
		SELECT u.id, u.name, u.email, u.token, u.start, u.extra_vacation, l.id, l.start, l.end, l.type, l.approved, l.user_id
		FROM users u
		INNER JOIN leaves l ON u.id = l.user_id
	`)

	if err != nil {
		return nil, err
	}

	// Map to store users by ID and collect leaves as child field for each user
	users := map[int]*model.User{}

	for rows.Next() {
		user := model.User{}
		leave := model.Leave{}

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

func (us *UserService) LeavesUsed(u *model.User, leaveTypes []string, periodStart, periodEnd *time.Time) int {
	var used int

	// gives ?,?,? for IN (?, ?, ?) if leaveTypes = []string{"vacation", "dayoff", "sick"}
	leaveTypesPlaceholders := strings.Repeat("?,", len(leaveTypes)-1) + "?"

	stmt := fmt.Sprintf(
		`SELECT id, start, end, type, approved, user_id
		FROM leaves
		WHERE user_id = ? AND start >= ? AND end <= ? AND type IN (%s)`, leaveTypesPlaceholders)

	args := []interface{}{u.ID, periodStart, periodEnd}
	for _, leaveType := range leaveTypes {
		args = append(args, leaveType)
	}

	rows, err := us.db.Query(stmt, args...)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		l := model.Leave{}
		err := rows.Scan(
			&l.ID,
			&l.Start,
			&l.End,
			&l.Type,
			&l.Approved,
			&l.UserID,
		)
		if err != nil {
			panic(err)
		}

		used += int(l.Duration().Hours() / 24)
	}
	rows.Close()

	return used
}
