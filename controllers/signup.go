package controllers

import (
	"net/http"
	"regexp"
	"strings"

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
	_, ok := formParams["fname"]
	if !ok || formParams["fname"][0] == "" {
		feedback = append(feedback, "You need to provide your First Name")
	}
	_, ok = formParams["sname"]
	if !ok || formParams["sname"][0] == "" {
		feedback = append(feedback, "You need to provide your Last Name")
	}
	// Check an eduroam value is submitted
	// If not then the user is signing up using a personal email
	if _, ok := formParams["eduroam"]; ok {
		eduroam := formParams["eduroam"][0]
		if eduroam == "" {
			feedback = append(feedback, "You need to provide your York Email")
		} else {
			// Ignore an added @york.ac.uk (since we assume it)
			eduroam = strings.TrimSuffix(eduroam, "@york.ac.uk")
			match, _ := regexp.MatchString("^[a-z]{1,6}[0-9]{1,6}$", eduroam)
			if !match {
				feedback = append(feedback, "The @york.ac.uk email you provided seems invalid")
			}
			formParams["eduroam"][0] = eduroam
		}
	} else {
		if _, ok = formParams["email"]; !ok {
			feedback = append(feedback, "You need to provide your email address")
		}
	}
	_, ok = formParams["phone"]
	if !ok || formParams["phone"][0] == "" {
		delete(formParams, "phone")
	}

	//If they are then post them off to the API
	if len(feedback) == 0 {
		sm := models.NewSignUpModel(gic.session)
		created, err := sm.Post(formParams)
		if err != nil {
			gic.handleError(w, r, err, "SignUpModel.Post")
			return
		}
		if !created {
			feedback = append(feedback, "Looks like you already have an account!")
			feedback = append(feedback, "Head over to "+gic.config.PageContext.FullURL+"/myradio to get started.")
			feedback = append(feedback, "If you can't sign in, click 'I've forgotten my login' to reset it.")
		}
	}

	data := struct {
		Feedback []string
	}{
		Feedback: feedback,
	}

	utils.RenderTemplate(w, gic.config.PageContext, data, "signedup.tmpl")
}
