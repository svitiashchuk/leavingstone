package tracker

import (
	"fmt"
	"leavingstone/internal/model"
	"leavingstone/internal/templ"
	"net/http"
	"strconv"
)

type TeamDetailsTemplateData struct {
	Team           *model.Team
	Members        []*model.MemberInfo
	WellbeingState string
	WellbeingIndex int
	*CommonTemplateData
}

func (app *App) TeamsHierarchy(w http.ResponseWriter, r *http.Request) {
	teams, err := app.teamService.AllTeams()
	if err != nil {
		app.internalError(w, err)
		return
	}

	w.Write([]byte(fmt.Sprintf("%+v", teams)))
}

func (app *App) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		team := &model.Team{
			Name: r.PostFormValue("name"),
		}
		err := app.teamService.CreateTeam(team)
		if err != nil {
			app.internalError(w, err)
			return
		}

		detailsPage := fmt.Sprintf("/teams/details?id=%d", team.ID)
		http.Redirect(w, r, detailsPage, http.StatusSeeOther)
		return
	}

	tmpl := app.templator.Page("team_create", nil)
	if err := tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r)); err != nil {
		app.internalError(w, err)
		return
	}
}

func (app *App) TeamDetails(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.internalError(w, err)
		return
	}

	team, err := app.teamService.FindByID(teamID)
	if err != nil {
		app.internalError(w, err)
		return
	}

	members, err := app.us.TeamMembers(teamID)
	if err != nil {
		app.internalError(w, err)
		return
	}

	tmpl := app.templator.Page("team_details", templ.Funcs())
	if err := tmpl.ExecuteTemplate(w, "layout", &TeamDetailsTemplateData{
		Team:               team,
		Members:            members,
		WellbeingState:     "good",
		WellbeingIndex:     87,
		CommonTemplateData: app.commonTemplateData(r),
	}); err != nil {
		app.internalError(w, err)
		return
	}
}
