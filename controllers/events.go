package controllers

import (
	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
	"log"
	"net/http"
	"strings"
)

// EventController is the controller for the ical service.
type EventController struct {
	Controller
}

// NewEventController returns a new EventController with the MyRadio
// session s and configuration context c.
func NewEventController(s *myradio.Session, c *structs.Config) *EventController {
	return &EventController{Controller{session: s, config: c}}
}

func (ec *EventController) Get(w http.ResponseWriter, r *http.Request) {
	em := models.NewEventsModel(ec.session)

	events, err := em.Get()
	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		Events []myradio.Event
	}{
		Events: events,
	}

	err = utils.RenderICal(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (ec *EventController) GetMeetings(w http.ResponseWriter, r *http.Request) {
	em := models.NewEventsModel(ec.session)

	events, err := em.Get()
	if err != nil {
		log.Println(err)
		return
	}

	var meetings []myradio.Event
	for _, event := range events {
		if strings.Contains(event.Title, "meeting") || strings.Contains(event.Title, "Meeting") {
			meetings = append(meetings, event)
		}
	}

	data := struct {
		Events []myradio.Event
	}{
		Events: meetings,
	}

	err = utils.RenderICal(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}
