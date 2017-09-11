package models

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TODO: this should have its own simplified CSS format, rather than the
// specific one provided. Otherwise new formats will require conversion or a
// different loading implementation.

type Location struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
	Country     string  `json:"country"`
}

type ByName []Location

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func ReadCityData(file io.Reader) (results []Location, err error) {
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip first header line

	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Split(line, "\t")

		// Translate region codes into human names, falling back on original value
		regionName := record[10]
		regionCode := record[8] + record[10]
		if name, found := REGION_CODES[regionCode]; found {
			regionName = name
		}

		location := Location{
			ID:          record[0],
			Name:        record[1],
			DisplayName: fmt.Sprintf("%s, %s, %s", record[1], regionName, record[8]),
			Country:     record[8],
		}
		if location.Lat, err = strconv.ParseFloat(record[4], 32); err != nil {
			return nil, err
		}
		if location.Long, err = strconv.ParseFloat(record[5], 32); err != nil {
			return nil, err
		}

		results = append(results, location)
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return results, nil
}

// Mapping FIPS region codes to provinces/states
var REGION_CODES = map[string]string{
	"CA01": "Alberta",
	"CA02": "British Columbia",
	"CA03": "Manitoba",
	"CA04": "New Brunswick",
	"CA05": "Newfoundland and Labrador",
	"CA07": "Nova Scotia",
	"CA08": "Ontario",
	"CA09": "Prince Edward Island",
	"CA10": "Quebec",
	"CA11": "Saskatchewan",
	"CA12": "Yukon",
	"CA13": "Northwest Territories",
	"CA14": "Nunavut",
}
