package db

import (
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"strconv"

	"go-whyye/pkg/quote"
	"go-whyye/pkg/user"
)

// CreateSchema creates the database schema.
func (d *Database) CreateSchema() error {
	const (
		usersTable = `
            CREATE TABLE users (
                id INTEGER PRIMARY KEY,
                name TEXT NOT NULL
            );
        `

		quotesTable = `
            CREATE TABLE quotes (
                user_id INTEGER NOT NULL REFERENCES users(id),
                message TEXT NOT NULL
            );
        `
	)

	_, err := d.db.Exec(usersTable)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(quotesTable)
	if err != nil {
		return err
	}

	return nil
}

// SeedDatabase reads the users_quotes.csv file and populates the database with users and quotes.
func (d *Database) SeedDatabase() error {
	const csvFile = "pkg/db/users_quotes.csv"

	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)

	// Read CSV header
	_, err = r.Read()
	if err != nil {
		return err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		id, _ := strconv.Atoi(record[0])
		user, quote := user.User{ID: id, Name: record[1]}, quote.Quote{UserID: id, Message: record[2]}

		_, err = d.db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.ID, user.Name)
		if err != nil {
			return err
		}
		_, err = d.db.Exec("INSERT INTO quotes (user_id, message) VALUES (?, ?)", user.ID, quote.Message)
		if err != nil {
			return err
		}

		fmt.Printf("Seeded user %d with quote '%s'\n", id, record[2])
	}
	return nil
}

// DropSchema drops the database schema.
func (d *Database) DropSchema() error {
	const (
		dropUsersTable  = "DROP TABLE users;"
		dropQuotesTable = "DROP TABLE quotes;"
	)

	_, err := d.db.Exec(dropUsersTable)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(dropQuotesTable)
	if err != nil {
		return err
	}
	return nil
}
