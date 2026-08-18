package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"lets-go-edwards/internal/models"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql" // only to call init() inside
)

// Dependency injection - to make "logger" global
// Next step: put all Handlers as Methods of this struct; so that "logger"/"snippets" are visible to them
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Setup Session Manager
	sessionManager := scs.New()
	// sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Store = mysqlstore.New(db) // Session Manager will use MySQL on Server to store Sessions
	sessionManager.Lifetime = 12 * time.Hour  // Session will expire after 12 hours
	// Make sure User will send cookies to us only when HTTPS is used
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Change default TLS config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		// CipherSuites: []uint16{
		// 	tls.TLS_AES_128_GCM_SHA256,
		// 	tls.TLS_AES_256_GCM_SHA384,
		// },
	}

	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Timeouts
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Limit Request Header size to 0.5MB - will be add 4096 bytes more by default
		MaxHeaderBytes: 524288,
	}

	logger.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
