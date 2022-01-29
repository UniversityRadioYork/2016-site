package controller

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

func (c *Controller) HandleTeam() http.HandlerFunc {
	v := view.View{
		MainTmpl: "team.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		teamM := models.NewTeamModel(c.Session)
		vars := mux.Vars(r)
		alias := vars["alias"]
		team, heads, assistants, officers, err := teamM.Get(alias)
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		data := struct {
			Team       myradio.Team
			Heads      []myradio.Officer
			Assistants []myradio.Officer
			Officers   []myradio.Officer
		}{
			Team:       team,
			Heads:      heads,
			Assistants: assistants,
			Officers:   officers,
		}
		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Controller) HandleTeams() http.HandlerFunc {
	v := view.View{
		MainTmpl: "teams.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		teamM := models.NewTeamModel(c.Session)
		teams, err := teamM.GetAll()
		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}
		data := struct {
			Teams []myradio.Team
		}{
			Teams: teams,
		}
		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
