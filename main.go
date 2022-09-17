package main

import (
	"fmt"
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

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	logger.Println("listening on ", listenTo)
	if err := http.ListenAndServe(listenTo, nil); err != nil {
		log.Fatal(err)
	}
}
