package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/mholt/binding"

	"backend_coding_challenge/models"
)

type SuggestionsController struct {
	locations *models.Trie
}

func NewSuggestionsController(locations *models.Trie) *SuggestionsController {
	return &SuggestionsController{locations: locations}
}

func (c *SuggestionsController) HandleSuggestions(res http.ResponseWriter, req *http.Request) {
	// Parse input from query string
	form := &SuggestionForm{Limit: 10}
	if err := binding.Bind(req, form); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("SuggestionsController: %#v", form)

	// Find locations matching the query prefix
	matches := c.locations.FindMatches(form.Query, form.Limit)

	// Initialize the algorithm used to score results
	var scorer models.Scorer
	scorer = models.NewRelativeLengthScorer(matches, form.Query)

	// Construct result objects from the locations and apply scores
	var results []models.Result
	for _, location := range matches {
		score := scorer.Score(location)

		results = append(results, models.Result{
			Name:  location.DisplayName,
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

type SuggestionForm struct {
	Query string  // Prefix to query locations
	Lat   float64 // Longitude for sorting results by distance (optional)
	Long  float64 // Latitude for sorting results by distance (optional)
	Limit int     // Limit to this many results in response (default 10)
}

// for auto-binding and validation with mholt/binding
func (form *SuggestionForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Query: binding.Field{
			Form:         "q",
			Required:     true,
			ErrorMessage: "query parameter 'q' is required",
		},
		&form.Lat:   "latitude",
		&form.Long:  "longitude",
		&form.Limit: "limit",
	}
}
