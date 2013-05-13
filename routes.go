// package ovdp-website renders the public-facing website for openvoicedata.org
package main

import (
	"flag"
	"fmt"
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

