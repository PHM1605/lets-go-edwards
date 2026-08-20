package models

import (
	"database/sql"
	"os"
	"testing"
)

// Establish an sql.DB connection for our test database
func newTestDB(t *testing.T) *sql.DB {
	// multiStatements=true: means we run many SQL statements in 1 db.Exec() call
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// read .sql file for setup
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	// execute .sql file
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	// Callback to clean up monitored by "t"
	// "t" will call the callback when all test that use newTestDB() has finished
	t.Cleanup(func() {
		defer db.Close()
		// run "teardown" script
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}
