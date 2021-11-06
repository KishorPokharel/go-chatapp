package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KishorPokharel/chatapp/pkg/models"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

type application struct {
	config        config
	logger        *log.Logger
	models        models.Models
	templateCache map[string]*template.Template
	sessions      *scs.SessionManager
	chatroom      *room
}

type config struct {
	dbdsn string
	port  int
	env   string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "port to listen to")
	flag.StringVar(&cfg.env, "env", "development", "(development | production)")
	flag.StringVar(&cfg.dbdsn, "dbdsn", os.Getenv("CHATAPP_DB_DSN"), "database url")
	flag.Parse()

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "authsess"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = cfg.env == "production"
	sessionManager.Store = redisstore.New(pool)

	app := &application{
		config:   cfg,
		logger:   log.New(os.Stdout, "", log.LstdFlags),
		sessions: sessionManager,
	}

	templateCache, err := newTemplateCache("templates")
	if err != nil {
		app.logger.Fatalln(err)
	}
	app.templateCache = templateCache

	db, err := app.openDB()
	if err != nil {
		app.logger.Fatalln(err)
	}
	app.logger.Println("database connection successful")
	app.models = models.New(db)

	app.chatroom = newRoom()
	go app.chatroom.run()

	log.Fatal(app.serve())
}

func (app *application) openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", app.config.dbdsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	return db, err
}
