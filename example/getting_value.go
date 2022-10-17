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
	if len(os.Args) < 2 {
		log.Fatalf("No key given in argument")
	}

	key := os.Args[1]
	log.Printf("Looking for key %s", key)
	resp, err := http.Get(fmt.Sprintf("http://localhost:18080/get/%s", key))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		log.Printf("Got response: %d", resp.StatusCode)
	} else {
		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("Got value: `%s`\n", string(buf.Bytes()))
	}
}
