package main

import (
	openhabcli "github.com/ebuyan/ohyandex/pkg/openhab"
	"github.com/joho/godotenv"
	"homekit/internal/homekit"
	"homekit/internal/openhab"
	"log"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalln("No .env.local file")
	}

	config := openhab.NewConfig()
	client := openhabcli.NewClient()
	repo := openhab.NewRepository(config, client)
	store := homekit.NewStore()
	factory := homekit.NewFactory(repo)
	bridge := homekit.NewBridge(repo, factory, store)

	err = bridge.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
