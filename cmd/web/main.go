package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/itsjustvita/bookings/pkg/config"
	"github.com/itsjustvita/bookings/pkg/handlers"
	"github.com/itsjustvita/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is our main functions
func main() {
	// change this to true when in production

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("cannot create template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = true
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatalln(err)
}
