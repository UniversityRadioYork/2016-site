package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
)

func (c *Controller) HandlePeople() http.HandlerFunc {
	v := view.View{
		MainTmpl: "people.tmpl",
		AddTmpls: []string{"elements/current_and_next.tmpl"},
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pm := models.NewPeopleModel(c.Session)

		vars := mux.Vars(r)

		id, _ := strconv.Atoi(vars["id"])

		user, officerships, credits, currentAndNext, err := pm.Get(id)

		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		if userHidden := len(credits) == 0 && len(officerships) == 0; userHidden {
			c.HandleNotFound(w, r)
			return
		}

		data := struct {
			User           *myradio.User
			Officerships   []myradio.Officership
			ShowCredits    []myradio.ShowMeta
			CurrentAndNext *myradio.CurrentAndNext
		}{
			User:           user,
			Officerships:   officerships,
			ShowCredits:    credits,
			CurrentAndNext: currentAndNext,
		}

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
