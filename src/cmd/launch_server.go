package main

import (
	"github.com/forgeronvirtuel/cache-server/src/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	if len(os.Args) < 1 {
		logger.Fatal("No listening address and port given in arguments")
	}

	listenTo := os.Args[1]

	for _, r := range routes.CreateRouteList(logger) {
		logger.Println("Creating handler for", r.Path)
		http.HandleFunc(r.Path, r.HandleHttp)
	}

	logger.Println("listening on ", listenTo)
	if err := http.ListenAndServe(listenTo, nil); err != nil {
		log.Fatal(err)
	}
}
