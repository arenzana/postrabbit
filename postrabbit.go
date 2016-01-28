package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// Config contains various config data populated from YAML
type Config struct {
	DbName         string
	DbUser         string
	DbPassword     string
	DbHost         string
	DbPort         string
	SslMode        string
	RabbitHost     string
	RabbitPort     string
	RabbitVHost    string
	RabbitQueue    string
	RabbitUser     string
	RabbitPassword string
}

var (
	app            = kingpin.New("postrabbit", "A PostgreSQL/RabbitMQ Example")
	setupcommand   = app.Command("setup", "setup the database for the example")
	runcommand     = app.Command("run", "run the listener")
	addcommand     = app.Command("add", "add a URL to the table")
	consumecommand = app.Command("consume", "consume data")
)

func main() {

	config := Config{}
	filebytes, err := ioutil.ReadFile("prcreds.yaml")
	if err != nil {
		log.Fatal("Failed to read creds")
	}
	err = yaml.Unmarshal(filebytes, &config)
	if err != nil {
		log.Fatal("Failed to parse creds", err)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case setupcommand.FullCommand():
		setup(config)
	case runcommand.FullCommand():
		run(config)
	case addcommand.FullCommand():
		add(config)
	case consumecommand.FullCommand():
		consume(config)

	}
}
