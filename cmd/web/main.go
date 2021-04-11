package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/config"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/handlers"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/render"
	"log"
	"net/http"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// TODO: change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
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

	fmt.Println(fmt.Sprintf("Starting application on port %s", PORT))

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
