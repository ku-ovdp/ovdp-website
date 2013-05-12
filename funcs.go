package main

import (
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"
)

// functions accessible to the templates
var funcMap = template.FuncMap{
	"stat": getStat,
}

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

func renderError(err error) template.HTML {
	return template.HTML(fmt.Sprintf("(error contacting api) <!-- %s -->", err))
}