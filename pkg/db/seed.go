package db

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// CreateSchema creates the database schema.
func (d *Database) CreateSchema() error {
	const (
		usersTable = `
            CREATE TABLE users (
                id INTEGER PRIMARY KEY,
                name TEXT NOT NULL,
								created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            );
        `

		quotesTable = `
            CREATE TABLE quotes (
								id INTEGER PRIMARY KEY,
                user_id INTEGER NOT NULL REFERENCES users(id),
                message TEXT NOT NULL,
								created_at DATETIME DEFAULT CURRENT_TIMESTAMP

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

// SeedDatabase reads the users_quotes.csv file and populates the database with users and quotes.
func (d *Database) SeedDatabase() error {
	const csvFile = "pkg/db/seed.csv"

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
		name := record[1]
		quotes := strings.Split(record[2], "|")
		now := time.Now()

		// Create user and insert into users table
		_, err = d.db.Exec("INSERT INTO users (id, name, created_at) VALUES (?, ?, ?)", id, name, now)
		if err != nil {
			return err
		}

		// Insert quotes into quotes table
		for _, quote := range quotes {
			_, err = d.db.Exec("INSERT INTO quotes (user_id, message, created_at) VALUES (?, ?, ?)", id, quote, now)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Seeded user %d with quote(s): %s\n", id, strings.Join(quotes, ",\n "))
	}

	return nil
}

func PrepareDb() error {
	db, err := NewDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("Dropping schema...")
	err = db.DropSchema()
	if err != nil {
		// return err
	}

	fmt.Println("Creating schema...")
	err = db.CreateSchema()
	if err != nil {
		return err
	}

	fmt.Println("Seeding...")
	err = db.SeedDatabase()
	if err != nil {
		return err
	}

	return nil
}
