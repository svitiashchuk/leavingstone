package sqlite

import (
	"database/sql"
	"leavingstone/internal/model"
)

type TeamService struct {
	db *sql.DB
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{
		db: db,
	}
}

func (s *TeamService) AllTeams() ([]*model.Team, error) {
	query := "SELECT id, name FROM teams"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := []*model.Team{}
	for rows.Next() {
		team := &model.Team{}
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (s *TeamService) FindByID(id int) (*model.Team, error) {
	query := "SELECT id, name FROM teams WHERE id = ?"
	row := s.db.QueryRow(query, id)

	team := &model.Team{}
	err := row.Scan(&team.ID, &team.Name)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) CreateTeam(team *model.Team) error {
	query := "INSERT INTO teams (name) VALUES (?)"
	result, err := s.db.Exec(query, team.Name)
	if err != nil {
		return err
	}

	teamID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	team.ID = int(teamID)
	return nil
}

func (s *TeamService) AssignLead(teamID int, leadID int) error {
	// Check if the lead is already a member of the team
	memberQuery := "SELECT COUNT(*) FROM team_members WHERE team_id = ? AND member_id = ?"
	var count int
	err := s.db.QueryRow(memberQuery, teamID, leadID).Scan(&count)
	if err != nil {
		return err
	}

	// If the lead is not a member, add them as a member
	if count == 0 {
		s.AddMember(teamID, leadID)
	}

	// Update the team's lead
	updateQuery := "UPDATE teams SET lead_id = ? WHERE id = ?"
	_, err = s.db.Exec(updateQuery, leadID, teamID)
	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) AddMember(teamID int, userID int) error {
	query := "UPDATE users SET team_id = ? WHERE id = ?"
	_, err := s.db.Exec(query, teamID, userID)
	return err
}

func (s *TeamService) RemoveMember(teamID int, memberID int) error {
	query := "UPDATE users SET team_id = NULL WHERE id = ?"
	_, err := s.db.Exec(query, memberID)
	return err
}
