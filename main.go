// package ovdp-website renders the public-facing website for openvoicedata.org
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// if flag is true templates are recompiled for every request (for development)
var recompile = flag.Bool("recompile", false, "recompile templates?")

// api base url
var apiEndpoint = flag.String("apiEndpoint", "http://api.openvoicedata.org/v1/", "api endpoint")

// map urls to templates
func handlers() {
	http.HandleFunc("/contact", renderTemplate("contact.html"))
	http.HandleFunc("/privacy", renderTemplate("privacy.html"))
	http.HandleFunc("/about", renderTemplate("about.html"))
	http.HandleFunc("/", renderTemplate("index.html"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
}

func main() {
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	handlers()
	fmt.Printf("attempting to listen on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// compiles `templates/<templateName>` along with templates/common/* and the functions
// defined in main.funcMap
func compileTemplate(templateName string) *template.Template {
	t := template.New("")

	// compile common templates and include funcMap (see funcs.go)
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