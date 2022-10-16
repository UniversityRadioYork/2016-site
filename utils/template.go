package utils

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gedex/inflector"

	"github.com/UniversityRadioYork/2016-site/structs"
	"github.com/UniversityRadioYork/myradio-go"
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
// RenderTemplate logs any error that occurred when rendering the template.
func RenderTemplate(w http.ResponseWriter, context structs.PageContext, data interface{}, mainTmpl string, addTmpls ...string) {
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
		"url":       func(s string) string { return PrefixURL(s, context.URLPrefix) },
		"html":      renderHTML,
		"stripHtml": StripHTML,
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
		"now": func() time.Time {
			return time.Now()
		},
		"subTime": func(aRaw, bRaw interface{}) (time.Duration, error) {
			var a, b time.Time
			var err error
			a, err = coerceTime(aRaw)
			if err != nil {
				return 0, err
			}
			b, err = coerceTime(bRaw)
			if err != nil {
				return 0, err
			}
			return a.Sub(b), nil
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
		logTemplateError("Error parsing templates %v: %v", tmpls, err)
	}

	err = t.Execute(w, td)
	if err != nil {
		logTemplateError("Error executing template %s: %v", mainTmpl, err)
	}
}

func logTemplateError(errMsg string, args ...interface{}) {
	msg := fmt.Sprintf(errMsg, args...)
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			msg = fmt.Sprintf("%s [at %s:%d (%s)]", msg, file, line, fn.Name())
		} else {
			msg = fmt.Sprintf("%s [at %s:%d]", msg, file, line)
		}
	}
	log.Println(msg)
}

// renderHTML takes some html as a string and returns a template.HTML
//
// Handles plain text gracefully.
func renderHTML(value interface{}) template.HTML {
	return template.HTML(fmt.Sprint(value))
}
