package controller

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
)

func (c *Controller) HandleGetInvolved() http.HandlerFunc {
	v := view.View{
		MainTmpl: "getinvolved.tmpl",
		Context:  c.Config.PageContext,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		gim := models.NewGetInvolvedModel(c.Session)

		colleges, numTeams, listTeamMap, faqs, err := gim.Get()

		if err != nil {
			log.Println(err)
			c.HandleNotFound(w, r)
			return
		}

		sort.Slice(colleges, func(i, j int) bool {
			left, right := colleges[i].CollegeName, colleges[j].CollegeName
			if strings.Contains(left, "N/A") {
				return true
			}
			if strings.Contains(left, "Unknown") && !strings.Contains(right, "N/A") {
				return true
			}
			return left < right
		})

		data := struct {
			Colleges    []myradio.College
			NumTeams    int
			ListTeamMap map[int]*myradio.Team
			FAQs        *models.FAQ
		}{
			Colleges:    colleges,
			NumTeams:    numTeams,
			ListTeamMap: listTeamMap,
			FAQs:        faqs,
		}

		err = v.Render(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
