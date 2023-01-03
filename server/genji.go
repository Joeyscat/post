package main

import (
	"github.com/genjidb/genji"
)

func initGenji(path string) (*genji.DB, error) {
	db, err := genji.Open(path)
	if err != nil {
		return nil, err
	}

	err = initTables(db)

	return db, err
}

func initTables(db *genji.DB) error {
	var err error
	err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS post (
	        id              TEXT    PRIMARY KEY,
	        userId          TEXT	NOT NULL,
	        content			TEXT,
	        medias         	ARRAY,
	        comments        ARRAY,
			time          	INT		NOT NULL
	    )
	`)
	if err != nil {
		return err
	}
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS post_user_idx ON post(userId)
	`)
	if err != nil {
		return err
	}

	err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS media (
	        id              TEXT    PRIMARY KEY,
	        userId          TEXT	NOT NULL,
	        type          	INT		NOT NULL,
	        url            	TEXT	NOT NULL,
	        comments        ARRAY,
	        posted        	BOOL,
			time          	INT		NOT NULL
	    )
	`)
	if err != nil {
		return err
	}
	err = db.Exec(`
		CREATE INDEX IF NOT EXISTS media_user_idx ON media(userId)
	`)
	if err != nil {
		return err
	}

	err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS user (
	        id              TEXT    PRIMARY KEY,
	        username        TEXT 	NOT NULL,
	        password        TEXT 	NOT NULL,
	        avatar        	TEXT 	NOT NULL,
			time          	INT		NOT NULL
	    )
	`)
	return err
}
