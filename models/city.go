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

type City struct {
	Id      string  // 5881791
	Name    string  // Abbotsford
	Lat     float64 // 49.05798
	Long    float64 // -122.25257
	Country string  // CA
}

func ReadCityData(file io.Reader) (results []City, err error) {
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip first header line

	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Split(line, "\t")

		city := City{
			Id:      record[0],
			Name:    record[1],
			Country: record[8],
		}
		if city.Lat, err = strconv.ParseFloat(record[4], 32); err != nil {
			return nil, err
		}
		if city.Long, err = strconv.ParseFloat(record[5], 32); err != nil {
			return nil, err
		}

		results = append(results, city)
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return results, nil
}
