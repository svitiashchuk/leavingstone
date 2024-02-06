package sqlite

import (
	"database/sql"
	"leavingstone/internal/model"
	"time"
)

const DBTimeFormat = "2006-01-02 15:04:05"

type LeaveService struct {
	db *sql.DB
}

// TODO pass DB
func NewLeaveService() (*LeaveService, error) {
	db, err := sql.Open("sqlite3", DSN)
	if err != nil {
		return nil, err
	}

	return &LeaveService{db}, nil
}

func (ls *LeaveService) List(from, to time.Time, limit int) ([]*model.Leave, error) {
	rows, err := ls.db.Query(`
		SELECT
		l.id, l.user_id, l.start, l.end, l.type, l.approved
		FROM leaves l
		INNER JOIN users u ON l.user_id = u.id
		WHERE l.start >= ? and l.end <= ?
	`,
		from.Format(DBTimeFormat),
		to.Format(DBTimeFormat),
	)

	if err != nil {
		return nil, err
	}

	ll := []*model.Leave{}
	for rows.Next() {
		leave := model.Leave{}
		err = rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.Start,
			&leave.End,
			&leave.Type,
			&leave.Approved,
		)

		if err != nil {
			return nil, err
		}

		ll = append(ll, &leave)
	}

	return ll, nil
}

func (ls *LeaveService) Upcoming(userID int) ([]*model.Leave, error) {
	rows, err := ls.db.Query(`
		SELECT
		l.id, l.user_id, l.start, l.end, l.type, l.approved
		FROM leaves l
		INNER JOIN users u ON l.user_id = u.id
		WHERE l.user_id = ?
		AND l.end >= ?
		ORDER BY l.start ASC
	`,
		userID,
		time.Now().Format(DBTimeFormat),
	)

	if err != nil {
		return nil, err
	}

	ll := []*model.Leave{}
	for rows.Next() {
		leave := model.Leave{}
		err = rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.Start,
			&leave.End,
			&leave.Type,
			&leave.Approved,
		)

		if err != nil {
			return nil, err
		}

		ll = append(ll, &leave)
	}

	return ll, nil
}

func (ls *LeaveService) AllUpcoming() ([]*model.Leave, error) {
	rows, err := ls.db.Query(`
		SELECT
		l.id, l.user_id, l.start, l.end, l.type, l.approved, u.id, u.name, u.email, u.token, u.start, u.extra_vacation
		FROM leaves l
		INNER JOIN users u ON l.user_id = u.id
		WHERE l.end >= ?
		ORDER BY l.start ASC
	`,
		time.Now().Format(DBTimeFormat),
	)

	if err != nil {
		return nil, err
	}

	users := map[int]*model.User{}
	ll := []*model.Leave{}
	for rows.Next() {
		u := model.User{}

		leave := model.Leave{}
		err = rows.Scan(
			&leave.ID,
			&leave.UserID,
			&leave.Start,
			&leave.End,
			&leave.Type,
			&leave.Approved,
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Token,
			&u.Started,
			&u.ExtraVacation,
		)

		if _, ok := users[u.ID]; !ok {
			users[u.ID] = &u
		}
		leave.User = users[u.ID]

		if err != nil {
			return nil, err
		}

		ll = append(ll, &leave)
	}

	return ll, nil
}

func (ls *LeaveService) Create(userID int, from, to time.Time, leaveType string) error {
	_, err := ls.db.Exec(
		`INSERT INTO leaves (start, end, user_id, type, approved) VALUES(?, ?, ?, ?, ?)`,
		from.Format(DBTimeFormat),
		to.Format(DBTimeFormat),
		userID,
		leaveType,
		false,
	)

	return err
}

func (ls *LeaveService) Approve(id int) error {
	_, err := ls.db.Exec(
		`UPDATE leaves SET approved = true WHERE id = ?`,
		id,
	)

	return err
}

func (ls *LeaveService) Reject(id int) error {
	_, err := ls.db.Exec(
		`UPDATE leaves SET approved = false WHERE id = ?`,
		id,
	)

	return err
}

func (ls *LeaveService) Delete(id int) error {
	_, err := ls.db.Exec(
		`DELETE FROM leaves WHERE id = ?`,
		id,
	)

	return err
}
