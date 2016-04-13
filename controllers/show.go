package controllers

import (
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/UniversityRadioYork/2016-site/models"
	"log"
	"github.com/cbroglie/mustache"
	"strconv"
)

type ShowController struct {
	Controller
}

func NewShowController(s *myradio.Session, c *structs.Config) *ShowController {
	return &ShowController{Controller{session:s, config:c}}
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
		log.Fatal(err)
		return
	}

	// Render Template

	pd := struct {
		Show    myradio.ShowMeta
		Seasons []myradio.Season
	}{
		Show: *show,
		Seasons: seasons,
	}

	td := struct {
		Globals structs.Globals
	}{
		Globals: structs.Globals{
			PageContext: sc.config.PageContext,
			PageData: pd,
		},
	}

	output, err := mustache.RenderFile("views/show.mustache", td)

	if err != nil {//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	w.Write([]byte(output))

}
