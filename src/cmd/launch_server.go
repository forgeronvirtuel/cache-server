package main

import (
	"errors"
	"fmt"
	"github.com/forgeronvirtuel/cache-server/src/routes"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	if len(os.Args) < 1 {
		logger.Fatal("No listening address and port given in arguments")
	}

	listenTo := os.Args[1]

	logger.Println("listening on ", listenTo)

	mux := http.NewServeMux()

	routes := routes.CreateRouteList(logger)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		for _, r := range routes {
			if strings.HasPrefix(request.URL.String(), r.Path) {
				r.HandleHttp(writer, request)
				return
			}
		}
		logger.Printf("Calling path `%s` not found.", request.URL)
		writer.WriteHeader(http.StatusNotFound)
	})

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
