package controllers

import (
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// SearchController is the controller for the search page.
type SearchController struct {
	Controller
}

// NewSearchController returns a new SearchController with the MyRadio session s
// and configuration context c.
func NewSearchController(s *myradio.Session, c *structs.Config) *SearchController {
	return &SearchController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the search page, writing to w.
func (sc *SearchController) Get(w http.ResponseWriter, r *http.Request) {
	// Check if they've landed or they've searched
	var term = r.URL.Query().Get("term")
	var searching = (term != "")
	var results []myradio.ShowMeta
	var err error

	if searching {
		// Contact the DB and get search results
		sm := models.NewSearchModel(sc.session)

		results, err = sm.Get(term)

		if err != nil {
			sc.handleError(w, r, err, "SearchModel.Get")
			return
		}
	}

	data := struct {
		Searching  bool
		Results    []myradio.ShowMeta
		NumResults int
		BaseURL    string
		Term       string
	}{
		Searching:  searching,
		Results:    results,
		NumResults: len(results),
		BaseURL:    r.URL.Path,
		Term:       term,
	}

	utils.RenderTemplate(w, sc.config.PageContext, data, "search.tmpl")
}
