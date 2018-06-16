package controllers

import (
	"log"
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
	var source = r.URL.Query().Get("source")

	var errorMsg string
	var searching = false

	if len(term) < 3 {
		errorMsg = "Your search term is too short. Try using more keywords."
	} else {
		searching = true
	}

	var showResults []myradio.ShowMeta
	var podcastResults []myradio.Podcast
	var peopleResults []myradio.UserSearch

	var err error

	if searching {
		// Contact the DB and get search results
		sm := models.NewSearchModel(sc.session)

		switch source {
		case "show":
			showResults, err = sm.GetShows(term)
		case "podcast":
			podcastResults, err = sm.GetPodcasts(term)
		case "people":
			peopleResults, err = sm.GetUsers(term)
		default:
			errorMsg = "You didn't select something to search for. Please select shows or podcasts."
		}

		if err != nil {
			log.Println(err)
			return
		}
	}
	var numResults = len(showResults) + len(podcastResults) + len(peopleResults)
	data := struct {
		Searching      bool
		Source         string
		ShowResults    []myradio.ShowMeta
		PodcastResults []myradio.Podcast
		PeopleResults  []myradio.UserSearch
		NumResults     int
		BaseURL        string
		Term           string
		ErrorMsg       string
	}{
		Searching:      searching,
		Source:         source,
		ShowResults:    showResults,
		PodcastResults: podcastResults,
		PeopleResults:  peopleResults,
		NumResults:     numResults,
		BaseURL:        r.URL.Path,
		Term:           term,
		ErrorMsg:       errorMsg,
	}

	err = utils.RenderTemplate(w, sc.config.PageContext, data, "search.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
