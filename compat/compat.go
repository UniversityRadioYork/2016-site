package compat

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
)

func OldConfig() *structs.Config {
	var config structs.Config
	toml.DecodeFile("config.toml", &config)
	if config.Schedule.StartHour != 0 {
		utils.StartHour = config.Schedule.StartHour
	}
	config.PageContext.CurrentYear = time.Now().Year()
	return &config
}
