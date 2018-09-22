package controllers

import (
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// ChartController is the controller for the URYPlayer Chart pages.
type ChartController struct {
	Controller
}

// NewChartController returns a new ChartController with the MyRadio session s
// and configuration context c.
func NewChartController(s *myradio.Session, c *structs.Config) *ChartController {
	return &ChartController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for a singular Chart page, writing to w.
func (ChartsC *ChartController) Get(w http.ResponseWriter, r *http.Request) {

	//Chartm := models.NewChartModel(ChartsC.session)

	//vars := mux.Vars(r)

	//id, _ := strconv.Atoi(vars["id"])

	//Chart, err := Chartm.Get(id)

	//if err != nil {
	//	log.Println(err)
	//	err = utils.RenderTemplate(w, ChartsC.config.PageContext, nil, "404.tmpl")
	//	return
	//}

	//if Chart.Status != "Published" {
	//	err = utils.RenderTemplate(w, ChartsC.config.PageContext, nil, "404.tmpl")
	//	return
	//}

	//data := struct {
	//Chart *myradio.Chart
	//}{
	//Chart: Chart,
	//}

	err := utils.RenderTemplate(w, ChartsC.config.PageContext, nil, "chart.tmpl")

	if err != nil {
		log.Println(err)
		return
	}

}
