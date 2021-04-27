package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/config"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/driver"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/handlers"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/helpers"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Starting application on port %s", PORT))

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// TODO: change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bedandbreakfast user=postgres password=usman123")
	if err != nil {
		log.Fatal("cannot connect to database. Shutting down...")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
