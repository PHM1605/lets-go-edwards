package models

import (
	"database/sql"
	"errors"
	"time"
)

// Create interface so that we can use both true DB wrapper AND mock DB wrapper
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert a Snippet into database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// what is the ID of the last row (i.e. from Snippet being inserted)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	// GET only the Snippet which isn't expired yet
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id=?`
	// Query 1 row only
	row := m.DB.QueryRow(stmt, id)
	// Prepare return Snippet
	var s Snippet
	// NOTE: number of parameters must match exactly number of columns in Statement "stmt"
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// not necessary; only to illustrate that we can use custom Error class
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	// All good
	return s, nil
}

// GET only the last 10 recent Snippets (which aren't expired yet)
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // NOTE: if "resultset" is not closed it will keep DB opens forever

	// Return value
	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// Still, we must check if "resultset" has any value after scanning
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// All good
	return snippets, nil
}
