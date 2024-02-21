package tracker

import (
	"fmt"
	"leavingstone/internal/model"
	"leavingstone/internal/templ"
	"net/http"
	"strconv"
	"text/template"
)

type TeamDetailsTemplateData struct {
	Team           *model.Team
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

	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/team_create.html",
	))

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
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

	team.Members, err = app.us.FindByTeamID(team.ID)
	if err != nil {
		app.internalError(w, err)
		return
	}

	tmpl := template.Must(
		template.
			New("team_details").
			Funcs(templ.Funcs()).
			ParseFiles(
				"frontend/src/templates/layout.html",
				"frontend/src/templates/team_details.html",
			),
	)
	if err := tmpl.ExecuteTemplate(w, "layout", &TeamDetailsTemplateData{
		Team:               team,
		WellbeingState:     "good",
		WellbeingIndex:     87,
		CommonTemplateData: app.commonTemplateData(r),
	}); err != nil {
		app.internalError(w, err)
	}
}
