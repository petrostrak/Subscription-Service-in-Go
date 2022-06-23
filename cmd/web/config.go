package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/petrostrak/Subscription-Service-in-Go/data"
)

type Config struct {
	Session   *scs.SessionManager
	DB        *sql.DB
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Wait      *sync.WaitGroup
	Models    data.Models
	Mailer    Mail
	ErrorChan chan error
	DoneChan  chan bool
}
