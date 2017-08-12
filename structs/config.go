package structs

// Config is a structure containing global website configuration.
//
// See the comments for Server and PageContext for more details.
type Config struct {
	Server      Server      `toml:"server"`
	PageContext PageContext `toml:"pageContext"`
}

// Server is a structure containing server configuration.
type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Timeout int    `toml:"timout"`
}

// PageContext is a structure containing static information to provide
// to all page templates.
//
// This contains the website's long and short names, as well as a directory
// of pages for navigation.
type PageContext struct {
	LongName        string `toml:"longName"`
	ShortName       string `toml:"shortName"`
	SiteDescription string `toml:"siteDescription"`
	URLPrefix       string `toml:"urlPrefix"`
	FullURL         string `toml:"fullURL"`
	MainTwitter     string `toml:"mainTwitter"`
	MainFacebook    string `toml:"mainFacebook"`
	NewsTwitter     string `toml:"newsTwitter"`
	Pages           []Page
}

// Page is a structure describing a page in the website navigation system.
type Page struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}
