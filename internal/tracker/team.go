package tracker

import (
	"fmt"
	"leavingstone/internal/model"
	"net/http"
	"strconv"
	"text/template"
)

type TeamDetailsTemplateData struct {
	team *model.Team
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

		http.Redirect(w, r, "/teams", http.StatusSeeOther)
	}

	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/create_team.html",
	))

	tmpl.ExecuteTemplate(w, "layout", app.commonTemplateData(r))
}

func (app *App) TeamDetails(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.internalError(w, err)
		return
	}

	team, err := app.teamService.FindTeamByID(teamID)
	if err != nil {
		app.internalError(w, err)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/team_details.html",
	))
	tmpl.ExecuteTemplate(w, "layout", &TeamDetailsTemplateData{
		team:               team,
		CommonTemplateData: app.commonTemplateData(r),
	})
}
