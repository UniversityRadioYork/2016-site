package controller

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
)

func (c *Controller) HandleSearch() http.HandlerFunc {
	v := view.View{
		MainTmpl: "search.tmpl",
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var term = r.URL.Query().Get("term")

		var results []myradio.ShowMeta
		var err error

		searching := term != ""
		if searching {
			sm := models.NewSearchModel(c.Session)

			results, err = sm.Get(term)

			if err != nil {
				log.Println(err)
				return
			}
		}

		data := struct {
			// maybe remove Searching in favour of checking Term directly in the template
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

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
