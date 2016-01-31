package controllers

import (
	"github.com/cbroglie/mustache"
	"github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/2016-site/structs"
	"log"
	"net/http"
)

type SearchController struct {
	Controller
}

func NewSearchController(s *myradio.Session, o *structs.Options) *SearchController {
	return &SearchController{Controller{session:s, options:o}}
}

func (ic *SearchController) Get(w http.ResponseWriter, r *http.Request) {

	// Check if they've landed or they've searched

	searching := false

	term := r.URL.Query().Get("term")

	searching = (term != "")

	if searching { // If searching

		// Contact the DB and get search results

		// Check if there are any results

		// If results

		// Show results

		// Else

		// "Sorry, no results"

		// End if

	} else {// Else "landed"
		// Do nothing
	}

	// Render Template

	td := struct {
		Globals   structs.Globals
		Searching bool
	}{
		Globals:    ic.options.Globals,
		Searching:  searching,
	}

	output, err := mustache.RenderFile("views/search.mustache", td)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(output))

}
