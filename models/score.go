package models

// A Scorer is used to calculate a score for each result returned by the server.
type Scorer interface {
	Score(Location) float64
}

// A RelativeLengthScorer scores results based on the length of their names
// relative to the length of the query. Longer names are given lower scores.
type RelativeLengthScorer struct {
	maxLength   int
	queryLength int
}

func NewRelativeLengthScorer(locations []Location, query string) *RelativeLengthScorer {
	// find maximum length in location names
	maxLength := 1
	for _, location := range locations {
		length := len(location.Name)
		if length > maxLength {
			maxLength = length
		}
	}

	return &RelativeLengthScorer{
		queryLength: len(query),
		maxLength:   maxLength,
	}
}

func (scorer *RelativeLengthScorer) Score(location Location) float64 {
	return 1.0 - (float64(len(location.Name)-scorer.queryLength) / float64(scorer.maxLength))
}
