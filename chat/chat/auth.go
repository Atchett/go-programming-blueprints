package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// implements the ServeHTTP method
// stores the http.Handler in the next field
type authHandler struct {
	next http.Handler
}

// satisfies the http.Handler interface
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check for an "auth" cookie
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth helper function creates authHandler that wraps
// any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles third parth login process
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 2 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Too few parameters")
		return
	}
	action := strings.ToLower(segs[2])
	provider := strings.ToLower(segs[3])
	switch action {
	case "login":
		log.Println("TODO: handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
