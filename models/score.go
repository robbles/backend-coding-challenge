package models

import "math"

// A Scorer is used to calculate a score for each result returned by the server.
type Scorer interface {
	Score(Location) float64
}

// A RelativeLengthScorer scores results based on the length of their names
// relative to the length of the query. Longer names are given lower scores.
type RelativeLengthScorer struct {
	queryLength int
}

func NewRelativeLengthScorer(locations []Location, query string) *RelativeLengthScorer {
	return &RelativeLengthScorer{
		queryLength: len(query),
	}
}

func (scorer *RelativeLengthScorer) Score(location Location) float64 {
	return inverseLengthScore(len(location.Name) - scorer.queryLength)
}

func inverseLengthScore(n int) float64 {
	return math.Exp2(-float64(n))
}
