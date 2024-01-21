package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dzivdzi/bookings/pkg/config"
	"github.com/dzivdzi/bookings/pkg/handlers"
	"github.com/dzivdzi/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Change this to true when in production
	app.InProduction = false

	// Create a session - initialize it and store it in a variable
	session = scs.New()
	// How long will the session persist - this means how long will the server know about the requester and associate items with it
	// This means that the server will REMEMBER the requester for 24h after the initial request (meaning that it will know the metadata)
	session.Lifetime = 24 * time.Hour
	// Should the cookie persist after the browser window is closed
	// If you want the session to dissapear once the browser is closed, set this to false
	session.Cookie.Persist = true
	// How strict you wanna be what site this cookie applies to - if you have multiple pages, sesion can be for only 1 page if you want
	session.Cookie.SameSite = http.SameSiteLaxMode
	// This will insist that the cookies are encrypted and that the connection is from https instead of http - in prod ALWAYS TRUE
	// as we connect to 8080, we don't hit a secure port
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
