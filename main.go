package main

import (
	"avoxi-api/server"
	"log"
)

func main() {
	var api server.Server
	err := api.Start(8080)
	if err != nil {
		log.Fatal(err)
	}
}