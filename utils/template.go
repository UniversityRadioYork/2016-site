package utils

import (
	"html/template"
	"net/http"
	"path/filepath"
	"log"
	"github.com/UniversityRadioYork/2016-site/structs"
)

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
func RenderTemplate(w http.ResponseWriter, context structs.PageContext, data interface{}, mainTmpl string, addTmpls ...string) {
	var err error = nil

	td := structs.Globals{
		PageContext: context,
		PageData:    data,
	}

	baseTmpls := append(
		[]string{
			"partials/header.tmpl",
			"partials/footer.tmpl",
			"elements/navbar.tmpl",
			"partials/base.tmpl",
			mainTmpl,
		},
		addTmpls...,
	)

	var tmpls []string
	for _, baseTmpl := range baseTmpls {
		tmpls = append(tmpls, filepath.Join("views", baseTmpl))
	}

	t := template.New("base.tmpl")
	t, err = t.ParseFiles(tmpls...)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, td)
	if err != nil {
		log.Println(err)
		return
	}
}
