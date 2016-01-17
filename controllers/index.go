package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/cbroglie/mustache"
	"log"
)

type IndexController struct {
	router *mux.Router
}

func NewIndexController(router *mux.Router) *IndexController {
	return &IndexController{router}
}

func (sc *IndexController) single(w http.ResponseWriter, r *http.Request) {

	// Call the model
	data := map[string]string{"content":"Hello World"}

	output, err := mustache.RenderFile("views/index.mustache", data)

	if (err != nil) {
		log.Fatal(err)
	} else {
		w.Write([]byte(output))
	}

}

func (sc *IndexController) Register(r string) {
	sc.router.HandleFunc(r, sc.single)
}
