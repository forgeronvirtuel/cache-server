package main

import (
	"errors"
	"fmt"
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

	logger.Println("listening on ", listenTo)

	mux := http.NewServeMux()

	for _, r := range routes.CreateRouteList(logger) {
		logger.Println("Creating handler for", r.Path)
		mux.HandleFunc(r.Path, (*r).HandleHttp)
	}

	server := http.Server{
		Addr:    listenTo,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
