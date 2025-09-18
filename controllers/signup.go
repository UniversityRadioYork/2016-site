package controllers

import (
	"log"
	"net/http"
	"regexp"
	"sort"
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

func (gic *SignUpController) Get(w http.ResponseWriter, r *http.Request) {
	gim := models.NewSignUpModel(gic.session)

	colleges, numTeams, listTeamMap, trainings, err := gim.Get()

	if err != nil {
		//@TODO: Do something proper here, render 404 or something
		log.Println(err)
		return
	}

	//Sort Colleges Alphabetically, with N/A and Unknown at the start
	sort.Sort(CollegeSorter(colleges))

	data := struct {
		Colleges    []myradio.College
		NumTeams    int
		ListTeamMap map[int]*myradio.Team
		Trainings   []models.SignUpTrainingSession
	}{
		Colleges:    colleges,
		NumTeams:    numTeams,
		ListTeamMap: listTeamMap,
		Trainings:   trainings,
	}

	err = utils.RenderTemplate(w, gic.config.PageContext, data, "signup.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
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
			match, _ := regexp.MatchString("^([a-z]|[A-Z]){1,6}[0-9]{1,6}$", eduroam)
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

	var trainingSignupResult int

	//If they are then post them off to the API
	if len(feedback) == 0 {
		sm := models.NewSignUpModel(gic.session)
		created, tsr, err := sm.Post(formParams)
		trainingSignupResult = tsr
		if err != nil {
			log.Println(err)
			feedback = append(feedback, "Oops. Something went wrong on our end.")
			feedback = append(feedback, "Please try again later")
		}
		if !created {
			feedback = append(feedback, "Looks like you already have an account!")
			feedback = append(feedback, "Head over to "+gic.config.PageContext.FullURL+"/myradio to get started.")
			feedback = append(feedback, "If you can't sign in, click 'I've forgotten my login' to reset it.")
		}
	}

	data := struct {
		Feedback             []string
		TrainingSignupResult int
	}{
		Feedback:             feedback,
		TrainingSignupResult: trainingSignupResult,
	}

	err := utils.RenderTemplate(w, gic.config.PageContext, data, "signedup.tmpl")

	if err != nil {
		log.Println(err)
		return
	}
}
