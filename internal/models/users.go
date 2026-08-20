package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Define Interface so that we can use either true DB wrapper OR mock DB wrapper
type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Interacting with DB
type UserModel struct {
	DB *sql.DB
}

// Get a User based on his ID
func (m *UserModel) Get(id int) (User, error) {
	var user User

	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`
	// Parse information from DB to our User variable
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}
	// all good
	return user, nil
}

// Insert new entry to "users" table
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) // 12 is 2^12 hashing iterations
	if err != nil {
		return err
	}

	// Insert to DB
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// try parsing "mySQLError"
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// "users_uc_email" is what we use when add constraint to MySQL table (check MySQL.md)
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	// all good
	return nil
}

// Verify if User exists; return "id" number if he does
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Buffer to get data from DB
	var id int
	var hashedPassword []byte

	// Query 1 User with his Email
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Verify if user-entered-password and hashed-password matches
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// all good => return User's ID
	return id, nil
}

// Verify if User exists with an ID
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)" // true or false
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

// Changing password
func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	// get current hashed password from DB
	var currentHashedPassword []byte
	stmt := "SELECT hashed_password FROM users WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}
	// compare password that Client gives & hashed password
	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}
	// if 2 current passwords match => create a new hashed password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	// update new hashed password in DB
	stmt = "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)

	return err
}
