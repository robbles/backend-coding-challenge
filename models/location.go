package models

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// TODO: this should have its own simplified CSS format, rather than the
// specific one provided. Otherwise new formats will require conversion or a
// different loading implementation.

type Location struct {
	Id      string
	Name    string
	Lat     float64
	Long    float64
	Country string
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

		location := Location{
			Id:      record[0],
			Name:    record[1],
			Country: record[8],
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
