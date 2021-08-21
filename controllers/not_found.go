package controllers

import (
	"fmt"
	"github.com/UniversityRadioYork/2016-site/models"
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

// NotFoundController is the controller for the 404 error page.
type NotFoundController struct {
	Controller
	shortURLs *models.ShortURLModel
}

// NewNotFoundController returns a new NotFoundController with the MyRadio
// session s and configuration context c.
func NewNotFoundController(c *structs.Config) *NotFoundController {
	shorts := models.NewShortURLsModel(c)
	go shorts.UpdateTimer()
	return &NotFoundController{
		Controller: Controller{config: c},
		shortURLs:  shorts,
	}
}

// Get handles the HTTP GET request r for the 404 page, writing to w.
// It first checks if a short URL matches, in which case it redirects to it.
func (sc *NotFoundController) Get(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path
	if slug[0] == '/' {
		slug = slug[1:]
	}
	if shortUrl := sc.shortURLs.Match(slug); shortUrl != nil {
		// Track the click asynchronously, for performance
		go func() {
			if err := sc.shortURLs.TrackClick(slug, utils.GetRequesterIP(sc.config, r)); err != nil {
				log.Println(fmt.Errorf("while tracking short URL click: %w", err))
			}
		}()
		http.Redirect(w, r, shortUrl.RedirectTo, http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(404)
	err := utils.RenderTemplate(w, sc.config.PageContext, nil, "404.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
