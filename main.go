package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/2016-site/utils"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/graceful"
)

func main() {
	log.SetFlags(log.Llongfile)

	config := &structs.Config{}
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Server.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.Server.SentryDSN,
			TracesSampleRate: 0.1,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("sentry.Init: %s", err))
		}
		defer sentry.Flush(2 * time.Second)
	}

	if config.Schedule.StartHour != 0 {
		utils.StartHour = config.Schedule.StartHour
	}

	config.PageContext.CurrentYear = time.Now().Year()

	s, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	log.Printf("Listening on: %s", l)
	graceful.Run(l, time.Duration(config.Server.Timeout), s)
}
