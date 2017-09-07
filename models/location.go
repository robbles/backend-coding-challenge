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

		// TODO: renderDisplayName function for translating state/country/province
		location := Location{
			ID:          record[0],
			Name:        record[1],
			DisplayName: fmt.Sprintf("%s, %s, %s", record[1], record[10], record[8]),
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
