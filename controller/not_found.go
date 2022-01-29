package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/UniversityRadioYork/2016-site/compat"
	"github.com/UniversityRadioYork/2016-site/models"

	"github.com/UniversityRadioYork/2016-site/utils"
)

func (c *Controller) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	c.notFoundView.Render(w, nil)
}

func (c *Controller) HandleShorts() http.HandlerFunc {
	// TODO: Remove compat patch
	shorts := models.NewShortURLsModel(compat.OldConfig(), c.Session)
	go shorts.UpdateTimer()

	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path
		if slug[0] == '/' {
			slug = slug[1:]
		}
		shortUrl := shorts.Match(slug)
		// would prefer an 'err != nil'
		if shortUrl == nil {
			c.HandleNotFound(w, r)
			return
		}
		go func() {
			// TODO: Remove compat patch
			reqIp, err := utils.GetRequesterIP(compat.OldConfig(), r)
			if err != nil {
				log.Println(fmt.Errorf("while getting requester IP: %w", err))
				return
			}
			err = shorts.TrackClick(
				shortUrl.ShortURLID,
				r.Header.Get("User-Agent"),
				reqIp,
			)
			if err != nil {
				log.Println(fmt.Errorf("while tracking short URL click: %w", err))
			}
		}()
		http.Redirect(w, r, shortUrl.RedirectTo, http.StatusTemporaryRedirect)
	}
}
