package controllers

import (
	"log"
	"net/http"
	"regexp"

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
func (gic *SignUpController) Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formParams := r.Form
	var feedback []string

	//Validate that necessary params are present and correct(enough)
	if formParams["fname"][0] == "" {
		feedback = append(feedback, "You need to provide your First Name")
	}
	if formParams["sname"][0] == "" {
		feedback = append(feedback, "You need to provide your Last Name")
	}
	if formParams["eduroam"][0] == "" {
		feedback = append(feedback, "You need to provide your York Email")
	} else {
		match, _ := regexp.MatchString("^[a-z]{1,6}[0-9]{1,6}$", formParams["eduroam"][0])
		if !match {
			feedback = append(feedback, "The @york.ac.uk email you provided seems invalid")
		}
	}
	if formParams["phone"][0] == "" {
		delete(formParams, "phone")
	}

	//If they are then post them off to the API
	if len(feedback) == 0 {
		sm := models.NewSignUpModel(gic.session)
		err := sm.Post(formParams)
		if err != nil {
			log.Println(err)
			feedback = append(feedback, "Oops. Something went wrong on our end.")
			feedback = append(feedback, "Please try again later")
		}
	}

	data := struct {
		Feedback []string
	}{
		Feedback: feedback,
	}

	err := utils.RenderTemplate(w, gic.config.PageContext, data, "signedup.tmpl")

	if err != nil {
		log.Println(err)
		return
	}
}
