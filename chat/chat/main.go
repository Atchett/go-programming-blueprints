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

const assetPath = "C:/Users/johns/Documents/go/src/bitbucket.org/johnpersonal/goblueprints/chat/chat/assets"

func main() {
	// get the addr flag, set to 8080 by default
	var addr = flag.String("addr", ":8080", "The address of the application.")
	flag.Parse() // parse the flags
	r := newRoom()
	// serve the assets
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(assetPath))))
	//r.tracer = trace.New(os.Stdout)
	// MustAuth triggers the authentication when user tries to access
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	// note the HandleFunc function to handle the loginHandler function
	// as loginHandler doesn't store any state, it's not an object
	http.HandleFunc("/auth/", loginHandler)
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
