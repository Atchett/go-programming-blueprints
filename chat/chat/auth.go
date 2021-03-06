package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
	// tidy up the output so we can see what we are getting
	outVal, _ := json.Marshal(segs)
	segsLength := len(segs)
	// log the number of degments
	//log.Println("Num segs : ", segsLength)
	if segsLength >= 4 {
		action := strings.ToLower(segs[2])
		provider := strings.ToLower(segs[3])
		switch action {
		case "login":
			// OLD - log to the console
			/* log.Println("TODO: handle login for", provider)
			log.Println("Segs", string(outVal)) */

			// get the provider that matches the object specified in the URL
			// e.g. google or facebook
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				// if there is an error write out with a non 200 code
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
				return
			}
			// get the location where we must send users
			// to start the authorization process
			// (nil, nil) arguments are for state and options - not in use for this app
			loginURL, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				// if there is an error write out with a non 200 code
				http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL %s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", loginURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		case "callback":
			provider, err := gomniauth.Provider(provider)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
				return
			}
			// lookup auth provider and call it's CompletAuth method
			creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			// relevant GetUser method from the provider
			user, err := provider.GetUser(creds)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
				return
			}
			authCookieValue := objx.New(map[string]interface{}{
				"name": user.Name(),
				// base encode the name field
				// ensure it won't contain any unpredictable characters
			}).MustBase64()
			// set the cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: authCookieValue,
				Path:  "/",
			})
			// redirect the user to the chat page
			w.Header().Set("Location", "/chat")
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			// write to the response (i.e. the page)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Auth action %s not supported.", action)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Too few params - %v", string(outVal))
	}

}
