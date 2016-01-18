package models

import (
	"encoding/json"
	"github.com/UniversityRadioYork/2016-site/structs"
	"log"
	"net/http"
)

type IndexModel struct {
	Model
}

// @TODO: Pass in the config options
func NewIndexModel() *IndexModel {
	return &IndexModel{}
}

func (m *IndexModel) Get() (data NowNextResponse, err error) {

	// @TODO: Move this into a config file!!!!
	url := "https://ury.org.uk/api/v2/timeslot/currentandnext/?api_key=9C4KCqywpDfzIk7OEhYO3tOjDJWftg2sZ65fKT5fTGCWvshnz5tinVt1MiqvETM4eZYDtQbRs13GoTCNB8HTsmQQlcDwFmRo8Xw3uHQoycYkumyTVGdXbxtt1S2Ow7RFbK"

	res, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return // @TODO: Is this the best way in Go??
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		return
	}

	return data, err

}

/**
 * @TODO: Look into whether or not this is the best way to do this
 */

type NowNextResponse struct {
	structs.Response
	Payload NowNextPayload `json:"payload"`
}

type NowNextPayload struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

type Show struct {
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Photo      string `json:"photo"`
	StartTime  uint64 `json:"start_time"`
	EndTime    uint64 `json:"end_time"`
	Presenters string `json:"presenters,omitempty"`
	Url        string `json:"url,omitempty"`
	Id         uint64 `json:"id,omitempty"`
}
