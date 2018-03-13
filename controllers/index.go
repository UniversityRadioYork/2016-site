package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// IndexController is the controller for the index page.
type IndexController struct {
	Controller
}

// RenderData is used to pass parameters into a common render function
type RenderData struct {
	CurrentAndNext *myradio.CurrentAndNext
	Banners        []myradio.Banner
	Teams          []myradio.Team
	MsgBoxError    bool
	JukeboxOnAir   bool
}

// NewIndexController returns a new IndexController with the MyRadio session s
// and configuration context c.
func NewIndexController(s *myradio.Session, c *structs.Config) *IndexController {
	return &IndexController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the index page, writing to w.
func (ic *IndexController) Get(w http.ResponseWriter, r *http.Request) {
	// This is where any form params would be parsed
	model := models.NewIndexModel(ic.session)

	currentAndNext, banners, teams, jukeboxOnAir, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}

	data := RenderData{
		CurrentAndNext: currentAndNext,
		Banners:        banners,
		Teams:          teams,
		JukeboxOnAir:   jukeboxOnAir,
		MsgBoxError:    false,
	}

	ic.render(w, data)
}

// Post handles the HTTP POST request r to the index page that sends a message. Writes to w.
func (ic *IndexController) Post(w http.ResponseWriter, r *http.Request) {
	// Parse message from request
	r.ParseForm()
	msg := r.Form.Get("message")

	// Get all the data for the webpage
	model := models.NewIndexModel(ic.session)

	currentAndNext, banners, teams, jukeboxOnAir, err := model.Get()

	if err != nil {
		log.Println(err)
		return
	}

	data := RenderData{
		CurrentAndNext: currentAndNext,
		Banners:        banners,
		Teams:          teams,
		JukeboxOnAir:   jukeboxOnAir,
		MsgBoxError:    false,
	}

	// Create the message model and send the message
	msgmodel := models.NewMessageModel(ic.session)
	err = msgmodel.Put(msg)
	if err != nil {
		// Set prompt if send fails
		data.MsgBoxError = true
	}

	ic.render(w, data)

}

func (ic *IndexController) render(w http.ResponseWriter, data RenderData) {
	// Render page
	err := utils.RenderTemplate(w, ic.config.PageContext, data, "index.tmpl", "elements/current_and_next.tmpl", "elements/banner.tmpl", "elements/message_box.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
