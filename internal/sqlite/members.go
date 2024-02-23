package sqlite

import (
	"leavingstone/internal/model"
	"time"
)

func (us *UserService) TeamMembers(teamID int) ([]*model.MemberInfo, error) {
	statsPeriodStart := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	statsPeriodEnd := time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC)

	rows, err := us.db.Query(`
		WITH
		approved_leaves AS (
			SELECT 
				users.id as user_id,
				CASE WHEN type = 'vacation' OR type = 'dayoff' THEN SUM(ROUND((julianday(leaves.end) - julianday(leaves.start)) + 0.5)) ELSE 0 END as vacations_total_days,
				CASE WHEN type = 'sick' THEN SUM(ROUND((julianday(leaves.end) - julianday(leaves.start)) + 0.5)) ELSE 0 END as sick_total_days,
				approved
			FROM leaves
			INNER JOIN users ON leaves.user_id = users.id
			WHERE approved = true and leaves.start >= ? and leaves.end < ? and users.team_id = ?
			GROUP BY users.id, approved
		), 
		today_status AS (
			SELECT 
				users.id as user_id,
				type
			FROM leaves INNER JOIN users ON leaves.user_id = users.id
			WHERE users.team_id = ? and date(leaves.start) <= date('now') and date('now') <= date(leaves.end)
			GROUP by users.id, type
		)
		SELECT 
			users.id as user_id,
			users.name,
			users.email,
			users.start,
			users.extra_vacation,
			COALESCE(vacations_total_days, 0),
			COALESCE(sick_total_days, 0),
			COALESCE(today_status.type, '')
		FROM users 
		LEFT JOIN approved_leaves ON users.id = approved_leaves.user_id
		LEFT JOIN today_status ON users.id = today_status.user_id
		WHERE team_id = ?
		GROUP BY users.id
		ORDER BY users.id;
	`, statsPeriodStart.Format(DBTimeFormat), statsPeriodEnd.Format(DBTimeFormat), teamID, teamID, teamID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]*model.MemberInfo, 0)
	for rows.Next() {
		member := &model.MemberInfo{}
		err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Email,
			&member.Started,
			&member.ExtraVacation,
			&member.VacationsUsed,
			&member.SickdaysUsed,
			&member.TodayStatus,
		)

		members = append(members, member)

		if err != nil {
			return nil, err
		}
	}

	return members, nil
}
