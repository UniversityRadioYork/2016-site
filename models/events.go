package models

import (
	"github.com/UniversityRadioYork/myradio-go"
)

type EventsModel struct {
	Model
}

func NewEventsModel(s *myradio.Session) *EventsModel {
	return &EventsModel{Model{session: s}}
}

func (m *EventsModel) Get() ([]myradio.Event, error) {
	//today := time.Now().Truncate(time.Hour * 24)
	//start := today.Format("20060102T150405Z")
	//end := today.AddDate(1,0,0).Format("20060102T150405Z")
	//return m.session.GetEventsInRange(start, end)

	return m.session.GetEventsNext(100)
}

func (m *EventsModel) GetMeetings([]myradio.Event) {

}
