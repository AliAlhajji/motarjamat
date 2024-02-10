package sqlitedb

import (
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

const filename string = "data.db"

var db *sqlx.DB

func Connect() (err error) {
	db, err = sqlx.Connect("sqlite3", filename)

	log.Println("Created database")

	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	for _, table := range createTables {

		_, err = db.Exec(table)
		if err != nil {
			return
		}
	}

	return nil
}

var createTables map[string]string = map[string]string{
	"users": `
		CREATE TABLE IF NOT EXISTS user (
	id INTEGER NOT NULL PRIMARY KEY,
	uuid TEXT,
	username TEXT,
	password TEXT DEFAULT "",
	email TEXT,
	name TEXT,
	join_date DATE,
	role text
	);
	`,

	"posts": `
			CREATE TABLE IF NOT EXISTS post(
				id INTEGER NOT NULL PRIMARY KEY,
				title TEXT,
				link TEXT,
				body TEXT,
				raw_body TEXT,
				user_id INTEGER,
				date DATE,
				FOREIGN KEY (user_id) REFERENCES user(id)
			);
	`,

	"categories": `
	CREATE TABLE IF NOT EXISTS category(
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT
		);
		`,

	"post_category": `
		CREATE TABLE IF NOT EXISTS post_category(
			post_id INTEGER,
			category_id INTEGET,
			FOREIGN KEY (post_id) REFERENCES user(id),
			FOREIGN KEY (category_id) REFERENCES user(id)
		);
	`,

	"vocabs": `
			CREATE TABLE IF NOT EXISTS vocab(
				id INTEGER NOT NULL PRIMARY KEY,
				en TEXT,
				ar TEXT,
				meaning TEXT
			);
	`,

	"post_vocab": `
	CREATE TABLE IF NOT EXISTS post_vocab (
		post_id INTEGER NOT NULL,
		vocab_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES post(id),
		FOREIGN KEY(vocab_id) REFERENCES vocab(id)
);
`,
}
