package controllers

import (
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/utils"
)

type SearchController struct {
	Controller
}

func NewSearchController(s *myradio.Session, c *structs.Config) *SearchController {
	return &SearchController{Controller{session:s, config:c}}
}

func (sc *SearchController) Get(w http.ResponseWriter, r *http.Request) {

	// Check if they've landed or they've searched

	var term string = r.URL.Query().Get("term")
	var searching bool = (term != "")
	var results []myradio.ShowMeta
	var err error

	if searching { // If searching

		// Contact the DB and get search results
		sm := models.NewSearchModel(sc.session)

		results, err = sm.Get(term)

		if err != nil {
			log.Println(err)
			return
		}

	}

	// Render Template

	data := struct {
		Searching  bool
		Results    []myradio.ShowMeta
		NumResults int
		BaseURL    string
		Term       string
	} {
		Searching:  searching,
		Results:    results,
		NumResults: len(results),
		BaseURL:    r.URL.Path,
		Term:        term,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "search.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
