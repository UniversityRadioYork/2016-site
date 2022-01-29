package controller

import (
	"log"
	"net/http"

	view "github.com/UniversityRadioYork/2016-site/views"
)

func (c *Controller) HandleStatic(tmpl string) http.HandlerFunc {
	v := view.View{
		MainTmpl: tmpl,
		Context:  c.Config.PageContext,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := v.Render(w, nil)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
