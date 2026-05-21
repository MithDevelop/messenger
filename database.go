package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase() error {

	var err error

	db, err = sql.Open(
		"sqlite",
		"messenger.db",
	)

	if err != nil {
		return err
	}
	_, err = db.Exec("PRAGMA journal_mode=WAL;")

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	return createTables()
}

func createTables() error {

	contactsTable := `
	CREATE TABLE IF NOT EXISTS contacts (
		peer_id TEXT PRIMARY KEY,
		username TEXT,
		added_at DATETIME,
		last_seen DATETIME,
		trusted BOOLEAN
	);
	`

	_, err := db.Exec(contactsTable)

	if err != nil {
		return err
	}

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		from_peer TEXT,
		to_peer TEXT,
		message_type TEXT,
		message TEXT,
		timestamp INTEGER
	);
	`

	_, err = db.Exec(messagesTable)

	if err != nil {
		return err
	}

	fmt.Println("Database initialized")

	return nil
}

func saveMessage(msg Message) error {

	query := `
	INSERT OR IGNORE INTO messages (
		id,
		from_peer,
		to_peer,
		message_type,
		message,
		timestamp
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		query,
		msg.ID,
		msg.From,
		msg.To,
		msg.Type,
		msg.Message,
		msg.Timestamp,
	)
	return err
}
