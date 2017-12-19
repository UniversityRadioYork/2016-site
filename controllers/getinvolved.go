package controllers

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// CollegeSorter exists so I can sort colleges properly
type CollegeSorter []myradio.College

// Implement sort.Interface
func (s CollegeSorter) Len() int {
	return len(s)
}
func (s CollegeSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s CollegeSorter) Less(i, j int) bool {
	if strings.Contains(s[i].CollegeName, "N/A") || strings.Contains(s[i].CollegeName, "Unknown") {
		return true
	}
	return s[i].CollegeName < s[j].CollegeName
}

// GetInvolvedController is the controller for the get involved page.
type GetInvolvedController struct {
	Controller
}

// NewGetInvolvedController returns a new GetInvolvedController with the MyRadio
// session s and configuration context c.
func NewGetInvolvedController(s *myradio.Session, c *structs.Config) *GetInvolvedController {
	return &GetInvolvedController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the get involved, writing to w.
func (gic *GetInvolvedController) Get(w http.ResponseWriter, r *http.Request) {

	gim := models.NewGetInvolvedModel(gic.session)

	colleges, numTeams, listTeamMap, err := gim.Get()

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	//Sort Colleges Alphabetically, with N/A and Unknown at the end
	sort.Sort(CollegeSorter(colleges))

	data := struct {
		Colleges    []myradio.College
		NumTeams    int
		ListTeamMap map[int]*myradio.Team
	}{
		Colleges:    colleges,
		NumTeams:    numTeams,
		ListTeamMap: listTeamMap,
	}

	err = utils.RenderTemplate(w, gic.config.PageContext, data, "getinvolved.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

}
