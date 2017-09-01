package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"backend_coding_challenge/controllers"
	"backend_coding_challenge/models"
)

func main() {
	var dataPath string
	var listenAddress string
	flag.StringVar(&dataPath, "data", "data/cities_canada-usa.tsv", "path to CSV source data")
	flag.StringVar(&listenAddress, "addr", ":8000", "TCP host:port to listen for requests on")
	flag.Parse()

	publicDir := "./public"

	f, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	// Read city data and insert into tree
	cities, err := models.ReadCityData(f)
	if err != nil {
		log.Fatal(err)
	}

	locations := models.NewTrie()

	for _, city := range cities {
		locations.Insert(city.Name, city)
	}

	suggestions := controllers.NewSuggestionsController(locations)

	mux := http.NewServeMux()
	mux.HandleFunc("/suggestions", suggestions.HandleSuggestions)
	mux.Handle("/", http.FileServer(http.Dir(publicDir)))

	log.Printf("Serving on %s...", listenAddress)
	log.Fatal(
		http.ListenAndServe(listenAddress, mux),
	)
}
