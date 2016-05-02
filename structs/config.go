package structs

type Config struct {
	Server      Server      `toml:"server"`
	PageContext PageContext `toml:"pageContext"`
}

type Server struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Timeout int    `toml:"timout"`
}

type PageContext struct {
	LongName  string `toml:"longName"`
	ShortName string `toml:"shortName"`
	Pages     []Page
}

type Page struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}
