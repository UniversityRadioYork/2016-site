package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/UniversityRadioYork/2016-site/structs"
	myradio "github.com/UniversityRadioYork/myradio-go"
	"github.com/gedex/inflector"
)

// TemplatePrefix is the constant containing the filepath prefix for templates.
const TemplatePrefix = "views"

// BaseTemplates is the array of 'base' templates used in each template render.
var BaseTemplates = []string{
	"partials/header.tmpl",
	"partials/footer.tmpl",
	"elements/navbar.tmpl",
	"partials/base.tmpl",
}

// RenderTemplate renders a 2016site template on the ResponseWriter w.
//
// This function automatically adds in the 2016site base templates, performs
// error handling, and builds a global context.
//
// The PageContext context gives the context for the page to be rendered, sent
// to the template as PageContext.
// The interface{} data gives the data to be sent to the template as PageData.
//
// The string mainTmpl gives the name, relative to views, of the main
// template to render.  The variadic argument addTmpls names any additional
// templates mainTmpl depends on.
//
// RenderTemplate returns any error that occurred when rendering the template.
func RenderTemplate(w http.ResponseWriter, context structs.PageContext, data interface{}, mainTmpl string, addTmpls ...string) error {
	var err error

	td := structs.Globals{
		PageContext: context,
		PageData:    data,
	}

	ownTmpls := append(addTmpls, mainTmpl)
	baseTmpls := append(BaseTemplates, ownTmpls...)

	var tmpls []string
	for _, baseTmpl := range baseTmpls {
		tmpls = append(tmpls, filepath.Join(TemplatePrefix, baseTmpl))
	}

	t := template.New("base.tmpl")
	t.Funcs(template.FuncMap{
		"url":  func(s string) string { return PrefixURL(s, context.URLPrefix) },
		"html": renderHTML,
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
		"showsToHours": func(shows []myradio.ShowMeta) int {
			//TODO: finish This
			return -5
		},
		"formatDuration": func(d time.Duration) string {
			s := d.String()
			var output = ""
			var startIndexLastNumerical = 0
			var suffix = ""
			for index := range s {
				if s[index] == []byte("d")[0] {
					suffix = " Day"
				} else if s[index] == []byte("h")[0] {
					suffix = " Hour"
				} else if s[index] == []byte("m")[0] {
					suffix = " Min"
				} else {
					suffix = ""
				}
				if suffix != "" {
					var value = ""
					var visible = true
					var plural = true
					value = string(s[startIndexLastNumerical:index])
					if len(value) == 1 {
						if value == "0" {
							visible = false
						} else if value == "1" {
							plural = false
						}
					}
					if visible {
						output = output + value + suffix
						if plural {
							output = output + "s"
						}
					}
					if index < len(s)-1 {
						startIndexLastNumerical = index + 1
					}
				}
			}
			return output
		},
		// TODO(CaptainHayashi): this is temporary
		"stripHTML": func(s string) string {
			d, err := StripHTML(s)
			if err != nil {
				return "Error stripping HTML"
			}
			return d
		},
		"week":   FormatWeekRelative,
		"plural": inflector.Pluralize,
	})
	t, err = t.ParseFiles(tmpls...)
	if err != nil {
		return err
	}

	return t.Execute(w, td)
}

// renderHTML takes some html as a string and returns a template.HTML
//
// Handles plain text gracefully.
func renderHTML(value interface{}) template.HTML {
	return template.HTML(fmt.Sprint(value))
}
