package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"backend_coding_challenge/models"
)

type SuggestionsController struct {
	locations *models.Trie
}

func NewSuggestionsController(locations *models.Trie) *SuggestionsController {
	return &SuggestionsController{locations: locations}
}

func (c *SuggestionsController) HandleSuggestions(res http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	limit := 10

	log.Printf("SuggestionsController: query=%#v", query)

	// Find locations matching the query prefix
	matches := c.locations.FindMatches(query, limit)

	// Initialize the algorithm used to score results
	var scorer models.Scorer
	scorer = models.NewRelativeLengthScorer(matches, query)

	// Construct result objects from the locations and apply scores
	var results []models.Result
	for _, location := range matches {
		score := scorer.Score(location)

		results = append(results, models.Result{
			Name:  fmt.Sprintf("%s, %s", location.Name, location.Country),
			Lat:   location.Lat,
			Long:  location.Long,
			Score: score,
		})
	}

	// Sort by score descending (should already be sorted by this point)
	sort.Sort(sort.Reverse(models.ResultsByScore(results)))

	// Write out the results
	if err := json.NewEncoder(res).Encode(results); err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, `{"error": "failed to marshal response as JSON"}`)
		return
	}
}
