package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Config contains various config data populated from YAML

func add(config Config) {

	purl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.SslMode)

	db, err := sql.Open("postgres", purl)

	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	var text = string(b)
	_, err = db.Exec(`INSERT INTO users(text) VALUES($1)`, text)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
