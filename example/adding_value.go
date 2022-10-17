package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("No key/value pair given in argument")
	}

	key := os.Args[1]
	val := os.Args[2]

	var buf io.Reader
	log.Printf("key = %s, value = %s", key, val)
	buf = bytes.NewReader([]byte(val))
	resp, err := http.Post(fmt.Sprintf("http://localhost:18080/add/%s", key), "application/json", buf)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Got response: %d", resp.StatusCode)
	} else {
		log.Printf("Value %s=%s added", key, val)
	}
}
