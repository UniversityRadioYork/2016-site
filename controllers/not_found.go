package controllers

import (
	"github.com/cbroglie/mustache"
	"log"
	"net/http"
	"github.com/UniversityRadioYork/2016-site/structs"
)

type NotFoundController struct {
	Controller
}

func NewNotFoundController(c *structs.Config) *NotFoundController {
	return &NotFoundController{Controller{config: c}}
}

func (sc *NotFoundController) Get(w http.ResponseWriter, r *http.Request) {

	output, err := mustache.RenderFile("views/404.mustache", map[string]string{})

	if err != nil {
		log.Println(err)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(output))
	}

}
