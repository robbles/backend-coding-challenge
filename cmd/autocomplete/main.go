package main

import (
	"backend_coding_challenge/models"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var dataPath string
	var limit int
	flag.StringVar(&dataPath, "data", "data/cities_canada-usa.tsv", "path to CSV source data")
	flag.IntVar(&limit, "limit", 10, "maximum number of results to return")
	flag.Parse()

	query := flag.Arg(0)

	f, err := os.Open(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	cities, err := models.ReadCityData(f)
	if err != nil {
		log.Fatal(err)
	}

	tree := models.NewRadixTree()

	for _, city := range cities {
		tree.Insert(city.Name)
	}

	for _, match := range tree.FindMatches(query, limit) {
		fmt.Println(match)
	}
}
