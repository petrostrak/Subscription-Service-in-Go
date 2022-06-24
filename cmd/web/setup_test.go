package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/petrostrak/Subscription-Service-in-Go/data"
)

var testApp Config

func TestMain(m *testing.M) {
	gob.Register(data.User{})

	// set up session
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	testApp = Config{
		Session:   session,
		DB:        nil,
		InfoLog:   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:  log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wait:      &sync.WaitGroup{},
		ErrorChan: make(chan error),
		DoneChan:  make(chan bool),
	}

	// create a dummy mailer
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	testApp.Mailer = Mail{
		Wait:       testApp.Wait,
		ErrorChan:  errorChan,
		MailerChan: mailerChan,
		DoneChan:   mailerDoneChan,
	}

	go func() {
		select {
		case <-testApp.Mailer.MailerChan:
		case <-testApp.Mailer.ErrorChan:
		case <-testApp.Mailer.DoneChan:
			return
		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErrorChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.DoneChan:
				return
			}
		}
	}()

	os.Exit(m.Run())
}