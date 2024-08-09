package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"go-whyye/pkg/quote"
	"go-whyye/pkg/user"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	app_env, exists := os.LookupEnv("APP_ENV")
	if !exists {
		app_env = "production"
	}

	db, err := sql.Open("sqlite3", "out/state/go-whyye."+app_env+".sql")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

// UsersListAll lists all users in the database.
func (d *Database) UsersListAll() ([]user.User, error) {
	const (
		getUsersQuery = `
            SELECT id, name FROM users;
        `
	)

	rows, err := d.db.Query(getUsersQuery)
	if err != nil {
		return nil, err
	}

	var users []user.User

	for rows.Next() {
		user := user.User{}
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UserById retrieves a user by their ID.
func (d *Database) UserById(id int) (*user.User, error) {
	const (
		getUserQuery = `
            SELECT id, name FROM users WHERE id = $1;
        `
	)

	row := d.db.QueryRow(getUserQuery, id)
	user := &user.User{}
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserByIdQuotes retrieves a user by their ID and returns their quotes.
func (d *Database) UserByIdQuotes(id int, limit int, offset int) ([]quote.Quote, error) {
	if id < 0 {
		return nil, errors.New("Invalid user id")
	}

	const (
		getUserQuotesQuery = `
            SELECT user_id, message FROM quotes WHERE user_id = $1 LIMIT $2 OFFSET $3;
        `
	)

	rows, err := d.db.Query(getUserQuotesQuery, id, limit, offset)
	if err != nil {
		return nil, err
	}

	var quotes []quote.Quote

	for rows.Next() {
		q := quote.Quote{}
		err := rows.Scan(&q.UserID, &q.Message)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, q)
	}

	return quotes, nil
}

func (d *Database) QuotesListAll(limit int, offset int) ([]quote.Quote, error) {
	const (
		getQuotesQuery = `
            SELECT user_id, message FROM quotes LIMIT $1 OFFSET $2;
        `
	)

	rows, err := d.db.Query(getQuotesQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	var quotes []quote.Quote

	for rows.Next() {
		q := quote.Quote{}
		err := rows.Scan(&q.UserID, &q.Message)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}

	return quotes, nil
}
