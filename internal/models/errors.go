package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// error when User logs in with invalid username and password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// error when User signs up with duplicated email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
