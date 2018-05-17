package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// compile the template once
	t.once.Do(func() {
		// parse the template in the templates folder
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// render itself using the data that can be extracted from the http.Request
	t.templ.Execute(w, r)
}

func main() {
	// get the addr flag, set to 8080 by default
	var addr = flag.String("addr", ":8080", "The address of the application.")
	flag.Parse() // parse the flags
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// get the room going as a go routine
	go r.run()
	// start the webserver using the reference to the flag addr value
	// call to flag.String returns type of *string
	// i.e. the address of a string variable where the value is stored
	// get the value, rather than the address, using *string
	// the pointer indirection operator
	log.Println("Starting the webserver on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
