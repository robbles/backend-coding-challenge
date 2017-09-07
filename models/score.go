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

func NewRelativeLengthScorer(query string) *RelativeLengthScorer {
	return &RelativeLengthScorer{
		queryLength: len(query),
	}
}

func (scorer *RelativeLengthScorer) Score(location Location) float64 {
	return InverseLengthScore(len(location.Name) - scorer.queryLength)
}

func InverseLengthScore(n int) float64 {
	return math.Exp2(-float64(n))
}

// A GeoDistanceScorer scores results based on their distance from the latitude
// and longitude provided in the query.
type GeoDistanceScorer struct {
	lat, long float64
}

func NewGeoDistanceScorer(lat, long float64) *GeoDistanceScorer {
	return &GeoDistanceScorer{
		lat: lat, long: long,
	}
}

// TODO: test with sample points and edge conditions
func (scorer *GeoDistanceScorer) Score(location Location) float64 {
	return DistanceScore(scorer.lat, scorer.long, location.Lat, location.Long)
}

// Calculate the distance between two points as a fraction of half the Earth's
// circumference (the maximum distance).
func CosineDistance(lat1, long1, lat2, long2 float64) float64 {
	lat1, long1, lat2, long2 = radians(lat1), radians(long1), radians(lat2), radians(long2)
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(long2-long1))

	// ensure rounding errors at small scales can never return NaN
	if math.IsNaN(dist) {
		return 0.0
	}

	// the maximum value for dist is pi, because this uses a sphere of unit radius
	return dist / math.Pi
}

func DistanceScore(lat1, long1, lat2, long2 float64) float64 {
	return 1.0 - CosineDistance(lat1, long1, lat2, long2)
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
