package structs

type Options struct {

		Server Server `toml:"server"`

}

type Server struct {

	Address	string	`toml:"address"`
	Port	int		`toml:"port"`
	Timeout int		`toml:"timout"`

}
