package controllers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// SubtypeController is the controller for the URY Show Subtype pages.
type SubtypeController struct {
	Controller
}

//NewSubtypeController returns a new controller for subtypes with:
// s: myradioSession
// c: structs config
func NewSubtypeController(s *myradio.Session, c *structs.Config) *SubtypeController {
	return &SubtypeController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the  page, writing to w.
func (subtypeCon *SubtypeController) Get(w http.ResponseWriter, r *http.Request) {
	subtypeM := models.NewSubtypeModel(subtypeCon.session)

	vars := mux.Vars(r)
	alias, _ := vars["alias"]

	subtype, err = subtypeM.Get(alias)
	if err != nil {
		log.Println(err)
		utils.RenderTemplate(w, subtypeCon.config.PageContext, nil, "404.tmpl")
		return
	}

	data := struct {
		Subtype  myradio.ShowSeasonSubtype
		Upcoming []myradio.Timeslot
		Recent   []myradio.Timeslot
	}{
		Subtype:  subtype,
		Upcoming: nil,
		Recent:   nil,
	}

	err := utils.RenderTemplate(w, subtypeCon.config.PageContext, data, "subtype.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
