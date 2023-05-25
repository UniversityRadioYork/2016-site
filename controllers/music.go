package controllers

import (
	"log"
	"net/http"
  "encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/UniversityRadioYork/myradio-go"
)

// MusicController is the controller for the Music page.
type MusicController struct {
	Controller
}

// Post represents an individual post in the RSS feed
type Post struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
  PubDate     string `xml:"pubDate"`
  Updated     string `xml:"updated"`
  Content     string `xml:"encoded"`
}

// NewMusicController returns a new MusicController with the MyRadio session s
// and configuration context c.
func NewMusicController(s *myradio.Session, c *structs.Config) *MusicController {
	return &MusicController{Controller{session: s, config: c}}
}

// Get handles the HTTP GET request r for the Music page, writing to w.
func (sc *MusicController) Get(w http.ResponseWriter, r *http.Request) {
  // URL of the RSS feed to query
  url := sc.config.Music.RSSFeed

  isError := false

  // Make HTTP GET request
  resp, err := http.Get(url)
  if err != nil {
    fmt.Printf("Failed to query URL: %v\n", err)
    isError = true
  }
  defer resp.Body.Close()

  // Read the response body
  xmlData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("Failed to read response body: %v\n", err)
    isError = true
  }

  // Parse the XML data into a struct
  var rss struct {
    Channel struct {
      Items []Post `xml:"item"`
    } `xml:"channel"`
  }

  err = xml.Unmarshal(xmlData, &rss)
  if err != nil {
    fmt.Printf("Failed to parse XML: %v\n", err)
    isError = true
  }

  // set content of each post to to scrape the first <em> tag
  for i := range rss.Channel.Items {
    var text string
    text, err = utils.ExtractFirstEmTagContent(string(rss.Channel.Items[i].Content))
    if err == nil {
      rss.Channel.Items[i].Content = text
    }
  }

  // If there was an error, render the error page
  if isError {
    err = utils.RenderTemplate(w, sc.config.PageContext, nil, "404.tmpl")
    if err != nil {
      log.Println(err)
      return
    }
  }

	err = utils.RenderTemplate(w, sc.config.PageContext, rss.Channel.Items, "music.tmpl")
	if err != nil {
		log.Println(err)
		return
	}
}
