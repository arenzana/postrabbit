package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Config contains various config data populated from YAML

func setup(config Config) {

	purl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.SslMode)
	db, err := sql.Open("postgres", purl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE users (text string primary key, timestamp timestamp without timezone);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE OR REPLACE FUNCTION notify_trigger() RETURNS trigger AS $$
BEGIN
	PERFORM pg_notify('usertrigger', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TRIGGER urlbefore BEFORE INSERT ON users
    FOR EACH ROW EXECUTE PROCEDURE notify_trigger();`)
	if err != nil {
		log.Fatal(err)
	}
}
