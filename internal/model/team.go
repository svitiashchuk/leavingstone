package model

type Team struct {
	ID      int
	Name    string
	Members []*User
	Lead    *User
}

type TeamService interface {
	AllTeams() ([]*Team, error)
	FindTeamByID(id int) (*Team, error)
	FindTeamByName(name string) (*Team, error)
}
