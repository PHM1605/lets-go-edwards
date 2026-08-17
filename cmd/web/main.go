package main

import (
	"database/sql"
	"flag"
	"html/template"
	"lets-go-edwards/internal/models"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql" // only to call init() inside
)

// Dependency injection - to make "logger" global
// Next step: put all Handlers as Methods of this struct; so that "logger"/"snippets" are visible to them
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

// dsn: Data Source Name (Connection String)
func openDB(dsn string) (*sql.DB, error) {
	// get a connections POOL; not get Connection to 'snippetbox' yet
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// actually create a Connection with a simple Ping
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address") // NOTE: this is a Pointer
	// Command line flag for MySQL connection string
	dsn := flag.String("dsn", "abc:1234@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Structure logger (instead of standard log)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// logger := slog.New(
	// 	slog.NewTextHandler(
	// 		os.Stdout, &slog.HandlerOptions{AddSource: true},
	// 	),
	// )

	// Database
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a Form Decoder
	formDecoder := form.NewDecoder()

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
