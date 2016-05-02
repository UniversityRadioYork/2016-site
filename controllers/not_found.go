package controllers

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"log"
	"net/http"
)

type NotFoundController struct {
	Controller
}

func NewNotFoundController(c *structs.Config) *NotFoundController {
	return &NotFoundController{Controller{config: c}}
}

func (sc *NotFoundController) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	err := utils.RenderTemplate(w, sc.config.PageContext, nil, "404.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
