package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/models"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// SignUpController is the controller for processing signup requests.
type SignUpController struct {
	Controller
}

// NewSignUpController returns a new SignUpController with the MyRadio
// session s and configuration context c.
func NewSignUpController(s *myradio.Session, c *structs.Config) *SignUpController {
	return &SignUpController{Controller{session: s, config: c}}
}

// Post handles the HTTP POST request r for the get involved, writing to w.
func (gic *SignUpController) Get(w http.ResponseWriter, r *http.Request) {

	formParams := r.URL.Query()
	log.Println(formParams)

	sm := models.NewSignUpModel(gic.session)

	feedback, err := sm.Post(formParams)

	log.Println(feedback)

	if err != nil {
		log.Println(err)
	}

	data := struct {
		Feedback []string
	}{
		Feedback: feedback,
	}

	err = utils.RenderTemplate(w, gic.config.PageContext, data, "signedup.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
