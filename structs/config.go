package structs

import "time"

// Config is a structure containing global website configuration.
//
// See the comments for Server and PageContext for more details.
type Config struct {
	Server         Server          `toml:"server"`
	PageContext    PageContext     `toml:"pageContext"`
	Schedule       ScheduleConfig  `toml:"schedule"`
	ShortURLs      ShortURLsConfig `toml:"shortUrls"`
	TrustedProxies []string        `toml:"trustedProxies"`
}

// Server is a structure containing server configuration.
type Server struct {
	Address    string `toml:"address"`
	Port       int    `toml:"port"`
	Timeout    int    `toml:"timeout"`
	MyRadioAPI string `toml:"myradio_api"`
}

// PageContext is a structure containing static information to provide
// to all page templates.
//
// This contains the website's long and short names, as well as a directory
// of pages for navigation.
type PageContext struct {
	LongName         string                `toml:"longName"`
	ShortName        string                `toml:"shortName"`
	SiteDescription  string                `toml:"siteDescription"`
	URLPrefix        string                `toml:"urlPrefix"`
	FullURL          string                `toml:"fullURL"`
	MainTwitter      string                `toml:"mainTwitter"`
	MainFacebook     string                `toml:"mainFacebook"`
	MainInstagram    string                `toml:"mainInstagram"`
	NewsTwitter      string                `toml:"newsTwitter"`
	MyRadioAPIKey    string                `toml:"publicMyRadioAPIKey"`
	ODName           string                `toml:"odName"`
	Christmas        bool                  `toml:"christmas"`
	AprilFools       bool                  `toml:"aprilFools"`
	CIN              bool                  `toml:"cin"`
	CINLivestreaming bool                  `toml:"cinLivestreaming"`
	CINAPI           string                `toml:"cinAPI"`
	CINHashtag       string                `toml:"cinHashtag"`
	CINLive          string                `toml:"cinLive"`
	IndexCountdown   *IndexCountdownConfig `toml:"indexCountdown"`
	CacheBuster      string                `toml:"cacheBuster"`
	Pages            []Page
	Youtube          youtube
	Gmaps            gmaps
	CurrentYear      int
}

// ScheduleConfig is a structure configuring the schedule views.
type ScheduleConfig struct {
	Sustainer SustainerConfig `toml:"sustainer"`
	StartHour int             `toml:"startHour"`
}

// SustainerConfig is a structure describing the sustainer show.
type SustainerConfig struct {
	Name string `toml:"name"`
	Desc string `toml:"desc"`
}

type IndexCountdownBackgroundType string

const (
	CountdownBackgroundVideo IndexCountdownBackgroundType = "video"
	CountdownBackgroundImage IndexCountdownBackgroundType = "image"
)

type IndexCountdownBackground struct {
	Type   IndexCountdownBackgroundType `toml:"type"`
	Source string                       `toml:"source"`
}

type IndexCountdownConfig struct {
	Enabled    bool                     `toml:"enabled"`
	Until      tomlTime                 `toml:"until"`
	Callout    string                   `toml:"callout"`
	LinkTo     string                   `toml:"linkTo"`
	Background IndexCountdownBackground `toml:"background"`
}

// Page is a structure describing a page in the website navigation system.
type Page struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

type youtube struct {
	APIKey             string `toml:"apiKey"`
	SessionsPlaylistID string `toml:"sessionsPlaylistID"`
	CINPlaylistID      string `toml:"cinPlaylistID"`
	ChannelURL         string `toml:"channelURL"`
}

type gmaps struct {
	APIKey string  `toml:"apiKey"`
	Lat    float32 `toml:"latitude"`
	Lng    float32 `toml:"longitude"`
}

// ShortURLsConfig is a structure configuring the short-urls subsystem.
type ShortURLsConfig struct {
	// UpdateInterval is how often the short URLs should be refreshed, in seconds.
	UpdateInterval uint `toml:"updateInterval"`
}

type tomlTime struct {
	time.Time
}

func (t *tomlTime) UnmarshalText(text []byte) error {
	var err error
	t.Time, err = time.Parse(time.RFC3339, string(text))
	return err
}
