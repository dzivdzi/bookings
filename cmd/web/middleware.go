package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf -> Adds CSRF protection to all post requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	// We need to create a cookie and make sure the token is available on a perPage basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, // Change to true on production
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad -> Loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	// LoadAndSave -> provides middleware which automatically loads and saves session
	// data for the current request, and communicates the session token to and from
	// the client as a cookie
	return session.LoadAndSave(next)
}
