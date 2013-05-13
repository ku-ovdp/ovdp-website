package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// functions accessible to the templates
var funcMap = template.FuncMap{
	"stat":         getStat,
	"current_year": getCurrentYear,
}

// get a stat from the ovdp api
// @todo this should do some caching
func getStat(stat string) template.HTML {
	if r, err := http.Get(*apiEndpoint + "stats"); err != nil {
		return renderError(err)
	} else {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		var stats map[string]int
		if err = decoder.Decode(&stats); err != nil {
			return renderError(err)
		}
		return template.HTML(fmt.Sprint(stats[stat]))
	}
}

func getCurrentYear() int {
	return time.Now().Year()
}

func renderError(err error) template.HTML {
	return template.HTML(fmt.Sprintf("(error contacting api) <!-- %s -->", err))
}


//////////// templating helpers

// compiles `templates/<templateName>` along with templates/common/* and the functions
// defined in main.funcMap
func compileTemplate(templateName string) *template.Template {
	t := template.New("")

	// compile common templates and include funcMap
	t = template.Must(t.Funcs(funcMap).
		ParseGlob("templates/common/*html"))

	// compile particular template along with common tempaltes and funcMap
	return template.Must(t.ParseFiles("templates/" + templateName))
}

// given a template name, returns an http.HandlerFunc that renders it
func renderTemplate(templateName string) http.HandlerFunc {
	templates := compileTemplate(templateName)
	context := getContext(templateName)
	return func(w http.ResponseWriter, r *http.Request) {
		if *recompile {
			templates = compileTemplate(templateName)
			context = getContext(templateName)
		}
		if err := templates.ExecuteTemplate(w, templateName, context); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getContext(templateName string) map[string]interface{} {
	return map[string]interface{}{
		"api_endpoint": *apiEndpoint,
	}
}
