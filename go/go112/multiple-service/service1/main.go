package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		service2URI := os.Getenv("SERVICE2_URI")
		resp, err := http.Get(service2URI)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			http.Error(w, "error occur", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		payload, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "Ok, get response: %s", payload)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
