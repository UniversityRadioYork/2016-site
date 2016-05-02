package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ShowController struct {
	Controller
}

func NewShowController(s *myradio.Session, c *structs.Config) *ShowController {
	return &ShowController{Controller{session: s, config: c}}
}

func (sc *ShowController) Get(w http.ResponseWriter, r *http.Request) {

	// Do the pagination!!

	// Call the DB for the things

	// Show the things

}

func (sc *ShowController) GetShow(w http.ResponseWriter, r *http.Request) {

	sm := models.NewShowModel(sc.session)

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	show, seasons, err := sm.GetShow(id)

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	// Render Template
	data := struct {
		Show    myradio.ShowMeta
		Seasons []myradio.Season
	}{
		Show:    *show,
		Seasons: seasons,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "show.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
