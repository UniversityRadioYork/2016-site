package view

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/UniversityRadioYork/2016-site/config"
	"github.com/UniversityRadioYork/2016-site/utils"
	myradio "github.com/UniversityRadioYork/myradio-go"
	"github.com/gedex/inflector"
)

var BaseTemplates = []string{
	"partials/base.tmpl",
	"partials/header.tmpl",
	"partials/footer.tmpl",
	"elements/navbar.tmpl",
}

type View struct {
	MainTmpl string
	AddTmpls []string
	Context  config.PageContext

	mux      sync.Mutex
	template *template.Template
}

//go:embed *.tmpl elements/*.tmpl partials/*.tmpl
var views embed.FS

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	err := v.checkParse()
	if err != nil {
		return err
	}

	td := struct {
		PageData    interface{}
		PageContext config.PageContext
	}{
		PageContext: v.Context,
		PageData:    data,
	}

	return v.template.Execute(w, td)
}

func (v *View) checkParse() error {
	v.mux.Lock()
	defer v.mux.Unlock()

	if v.template != nil {
		return nil
	}
	if err := v.parse(); err != nil {
		return err
	}
	return nil
}

func (v *View) parse() error {
	ownTmpls := append(v.AddTmpls, v.MainTmpl)
	tmpls := append(BaseTemplates, ownTmpls...)

	base := BaseTemplates[0]
	name := base[strings.LastIndex(base, "/")+1:]
	t := template.New(name)
	t.Funcs(template.FuncMap{
		"url": func(s string) string { return utils.PrefixURL(s, v.Context.URLPrefix) },
		"html": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"stripHtml": utils.StripHTML,
		//Takes a splice of show meta and returns the last x elements
		"getLastShowMeta": func(a []myradio.ShowMeta, amount int) []myradio.ShowMeta {
			if len(a) < amount {
				return a
			}
			return a[len(a)-amount:]

		},
		//Takes a splice of seasons and returns the total number of episodes
		"showCount": func(seasons []myradio.Season) int {
			var c = 0
			for _, season := range seasons {
				//Something about JSON being read as a float 64 so needing to convert to an int
				c += int(season.NumEpisodes.Value.(float64))
			}
			return c
		},
		"formatDuration": func(d time.Duration) string {
			days := int64(d.Hours()) / 24
			hours := int64(d.Hours()) % 24
			minutes := int64(d.Minutes()) % 60
			seconds := int64(d.Seconds()) % 60

			segments := []struct {
				name  string
				value int64
			}{
				{"Day", days},
				{"Hour", hours},
				{"Min", minutes},
				{"Sec", seconds},
			}

			parts := []string{}

			for _, s := range segments {
				if s.value == 0 {
					continue
				}
				plural := ""
				if s.value != 1 {
					plural = "s"
				}

				parts = append(parts, fmt.Sprintf("%d %s%s", s.value, s.name, plural))
			}
			return strings.Join(parts, " ")
		},
		"formatTime": func(fmt string, t time.Time) string {
			return t.Format(fmt)
		},
		// TODO(CaptainHayashi): this is temporary
		"stripHTML": func(s string) string {
			d, err := utils.StripHTML(s)
			if err != nil {
				return "Error stripping HTML"
			}
			return d
		},
		"week":   utils.FormatWeekRelative,
		"plural": inflector.Pluralize,
	})
	_, err := t.ParseFS(views, tmpls...)
	if err != nil {
		return fmt.Errorf("unable to parse view: %v", err)
	}
	v.template = t
	return nil
}
