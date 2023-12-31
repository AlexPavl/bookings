package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlexPavl/bookings/cmd/pkg/config"
	"github.com/AlexPavl/bookings/cmd/pkg/handlers"
	"github.com/AlexPavl/bookings/cmd/pkg/render"

	"github.com/alexedwards/scs/v2"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true                  // for continuing session after exit
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict cookies
	session.Cookie.Secure = app.InProduction       // if secure connection (https - is secured)
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
