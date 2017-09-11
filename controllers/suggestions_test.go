package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"backend_coding_challenge/models"
)

// alias NewResult for quick shorthand
var result func(models.Location, float64) models.Result = models.NewResult

func TestSuggestionsController_HandleSuggestions(t *testing.T) {
	// sample locations
	victoria := models.Location{ID: "6174041", Name: "Victoria", DisplayName: "Victoria, 02, CA", Lat: 48.43294143676758, Long: -123.36930084228516, Country: "CA"}
	vista := models.Location{ID: "5406602", Name: "Vista", DisplayName: "Vista, CA, US", Lat: 33.20003890991211, Long: -117.24253845214844, Country: "US"}

	locations := models.NewTrie()
	locations.Insert("Victoria", victoria)
	locations.Insert("Vista", vista)

	suggestions := NewSuggestionsController(locations)

	tests := map[string]struct {
		query   string
		status  int
		results []models.Result
	}{
		"no query parameter": {
			"",
			400,
			nil,
		},
		"bad limit": {
			"q=Vi&limit=hello",
			400,
			nil,
		},
		"successful query with results": {
			"q=Vi",
			200,
			[]models.Result{
				result(vista, models.InverseLengthScore(3)),
				result(victoria, models.InverseLengthScore(6)),
			},
		},
		"successful query with results and limit": {
			"q=Vi&limit=1",
			200,
			[]models.Result{
				result(vista, models.InverseLengthScore(3)),
			},
		},
		"successful query without results": {
			"q=Nope",
			200,
			[]models.Result{},
		},
		"successful query with lat/long": {
			"q=Vi&latitude=48.43&longitude=-123.33",
			200,
			[]models.Result{
				result(victoria, models.DistanceScore(48.43, -123.33, victoria.Lat, victoria.Long)),
				result(vista, models.DistanceScore(48.43, -123.33, vista.Lat, vista.Long)),
			},
		},
		"query with lat/long does not limit before scoring/sorting": {
			"q=Vi&latitude=48.43&longitude=-123.33&limit=1",
			200,
			[]models.Result{
				result(victoria, models.DistanceScore(48.43, -123.33, victoria.Lat, victoria.Long)),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/suggestions?"+tt.query, nil)
			res := httptest.NewRecorder()

			suggestions.HandleSuggestions(res, req)

			if res.Code != tt.status {
				t.Fatalf("Unexpected HTTP status %v != %v", res.Code, tt.status)
			}

			if res.Code != 200 {
				return
			}

			results := []models.Result{}
			if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(results, tt.results) {
				t.Errorf("%#v != %#v", results, tt.results)
			}
		})
	}
}
