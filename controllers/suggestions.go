package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	results := c.locations.FindMatches(query, limit)

	if err := json.NewEncoder(res).Encode(results); err != nil {
		res.WriteHeader(500)
		fmt.Fprint(res, `{"error": "failed to marshal response as JSON"}`)
		return
	}
}
