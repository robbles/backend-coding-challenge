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

	// Initialize the algorithm used to score results
	var scorer models.Scorer
	var matches []models.Location
	var matchLimit int

	if form.Lat != nil && form.Long != nil {
		// Use geo distance for scoring when latitude and longitude are passed
		scorer = models.NewGeoDistanceScorer(*form.Lat, *form.Long)

		// Don't limit results before scoring, or close results may be excluded
		matchLimit = 0
	} else {
		// Fall back to scoring by length relative to the prefix otherwise
		scorer = models.NewRelativeLengthScorer(form.Query)
		matchLimit = form.Limit
	}

	matches = c.locations.FindMatches(form.Query, matchLimit)
	log.Printf("%d matches found for prefix query", len(matches))

	// Construct result objects from the locations and apply scores
	results := []models.Result{}
	for _, location := range matches {
		score := scorer.Score(location)
		results = append(results, models.NewResult(location, score))
	}

	// Sort by score descending (should already be sorted by this point)
	sort.Sort(sort.Reverse(models.ResultsByScore(results)))

	// Trim array of results to <limit>
	if form.Limit > 0 && len(results) > form.Limit {
		results = results[:form.Limit]
	}

	// Write out the results
	if err := json.NewEncoder(res).Encode(results); err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, `{"error": "failed to marshal response as JSON"}`)
		return
	}
}

type SuggestionForm struct {
	Query string   // Prefix to query locations
	Lat   *float64 // Longitude for sorting results by distance (optional)
	Long  *float64 // Latitude for sorting results by distance (optional)
	Limit int      // Limit to this many results in response (default 10)
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
