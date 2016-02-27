package web

import (
	"github.com/UniversityRadioYork/2016-site/controllers"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/UniversityRadioYork/myradio-go"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(o *structs.Options) (*Server, error) {

	s := Server{negroni.Classic()}

	session, err := myradio.NewSessionFromKeyFile()

	if err != nil {
		return &s, err;
	}

	router := mux.NewRouter()

	getRouter := router.Methods("GET").Subrouter()

	// Routes go in here
	nfc := controllers.NewNotFoundController(o)
	router.NotFoundHandler = http.HandlerFunc(nfc.Get)

	ic := controllers.NewIndexController(session, o)
	getRouter.HandleFunc("/", ic.Get)

	sc := controllers.NewSearchController(session, o)
	getRouter.HandleFunc("/search", sc.Get)

	showC := controllers.NewShowController(session, o)
//	getRouter.HandleFunc("/schedule/shows", showC.Get) // @TODO: Implement this
	getRouter.HandleFunc("/schedule/shows/{id:[0-9]+}", showC.GetShow)

	// End routes

	s.UseHandler(router)

	return &s, nil;

}
