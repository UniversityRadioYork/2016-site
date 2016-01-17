package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/cbroglie/mustache"
	"log"
)

type NotFoundController struct {
	router *mux.Router
}

func NewNotFoundController(router *mux.Router) *NotFoundController {
	return &NotFoundController{router}
}

func (sc *NotFoundController) single(w http.ResponseWriter, r *http.Request) {

	output, err := mustache.RenderFile("views/404.mustache", map[string]string{"content":"hello world"})

	if (err != nil) {
		log.Fatal(err)
	} else {
		w.Write([]byte(output))
	}

}

func (sc *NotFoundController) Register() {
	sc.router.NotFoundHandler = http.HandlerFunc(sc.single)
}
