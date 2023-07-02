package sqlite

import (
	"database/sql"
	"ptocker"
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

func (ls *LeaveService) List(from, to time.Time, limit int) ([]*ptocker.Leave, error) {
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

	ll := []*ptocker.Leave{}
	for rows.Next() {
		leave := ptocker.Leave{}
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

func (ls *LeaveService) Create(userID int, from, to time.Time, leaveType string) error {
	_, err := ls.db.Query(
		`INSERT INTO leaves (start, end, type, user_id) VALUES(?, ?, ?, ?)`,
		from.Format(DBTimeFormat),
		to.Format(DBTimeFormat),
		userID,
		leaveType,
	)

	return err
}

func (ls *LeaveService) Approve(id int) error {
	_, err := ls.db.Query(
		`UPDATE leaves SET approved = true WHERE id = ?`,
		id,
	)

	return err
}

func (ls *LeaveService) Reject(id int) error {
	_, err := ls.db.Query(
		`UPDATE leaves SET approved = true WHERE id = ?`,
		id,
	)

	return err
}

func (ls *LeaveService) Delete(id int) error {
	_, err := ls.db.Query(
		`DELETE FROM leaves WHERE id = ?`,
		id,
	)

	return err
}
