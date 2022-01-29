package controller

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	view "github.com/UniversityRadioYork/2016-site/views"
	"github.com/UniversityRadioYork/myradio-go"
)

type RenderData struct {
	CurrentAndNext *myradio.CurrentAndNext
	Banners        []myradio.Banner
	Teams          []myradio.Team
	Podcasts       []myradio.Podcast
	MsgBoxError    bool
	ShowOnAir      bool
}

func (c *Controller) HandleIndex() http.HandlerFunc {
	v := view.View{
		MainTmpl: "index.tmpl",
		AddTmpls: []string{"elements/current_and_next.tmpl", "elements/banner.tmpl", "elements/message_box.tmpl", "elements/index_countdown.tmpl"},
		Context:  c.Config.PageContext,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		model := models.NewIndexModel(c.Session)

		currentAndNext, banners, teams, podcasts, showOnAir, err := model.Get()

		if err != nil {
			log.Println(err)
			return
		}

		data := RenderData{
			CurrentAndNext: currentAndNext,
			Banners:        banners,
			Teams:          teams,
			Podcasts:       podcasts,
			ShowOnAir:      showOnAir,
			MsgBoxError:    false,
		}

		v.Render(w, data)
	}
}

func (c *Controller) HandleIndexPost() http.HandlerFunc {
	v := view.View{
		MainTmpl: "index.tmpl",
		AddTmpls: []string{"elements/current_and_next.tmpl", "elements/banner.tmpl", "elements/message_box.tmpl", "elements/index_countdown.tmpl"},
		Context:  c.Config.PageContext,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		msg := r.Form.Get("message")

		model := models.NewIndexModel(c.Session)

		currentAndNext, banners, teams, podcasts, showOnAir, err := model.Get()

		if err != nil {
			log.Println(err)
			return
		}

		data := RenderData{
			CurrentAndNext: currentAndNext,
			Banners:        banners,
			Teams:          teams,
			Podcasts:       podcasts,
			ShowOnAir:      showOnAir,
			MsgBoxError:    false,
		}

		msgmodel := models.NewMessageModel(c.Session)
		err = msgmodel.Put(msg)
		if err != nil {
			// Set prompt if send fails
			data.MsgBoxError = true
		}

		v.Render(w, data)
	}
}
