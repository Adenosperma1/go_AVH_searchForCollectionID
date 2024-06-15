package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Occurrence struct {
	CatalogNumber string `json:"raw_catalogNumber"`
	InstitutionCode string `json:"raw_institutionCode"`
	Genus string `json:"genus"`
	Species string `json:"species"` //Bug??? RETURNS Genus and species!!!
	StateProvince string `json:"stateProvince"`
	Latitude float32 `json:"decimalLatitude"`
	Longitude float32 `json:"decimalLongitude"`
	Year int `json:"year"`
	MonthNumber string `json:"month"`
	TimeStamp int64 `json:"eventDate"`
	//Day string `json:"day"` //NOT SUPPLIED IN JSON???
	Collectors []string `json:"collectors"`
	CollectorsNumber string `json:"recordNumber"`
	//"locality missing???"
}

type Response struct {
	TotalRecords int    `json:"totalRecords"`
	Query        string `json:"queryTitle"`
	Occurrences []Occurrence `json:"occurrences"`
}

func main() {
	//TODO read a csv file with id numbers

	collectionIds := []string{"CANB220604", "cbg8703783"}

	urlStart := "https://biocache-ws.ala.org.au/ws/occurrences/search?q=text%3A%22"
	urlEnd := "%22&disableAllQualityFilters=true&qc=data_hub_uid%3Adh9"

	searchResponses := searchForCollections(collectionIds, urlStart, urlEnd)

	// fmt.Println("All Results:")
	for _, response := range searchResponses {
		// fmt.Printf("Response Body %d: %s\n", i, string(response))

		// Create a struct obj to hold the decoded JSON data
		responseJSON := createJson(response)

		if responseJSON.TotalRecords == 0 {
			fmt.Println("ERROR: Not found on AVH:", getIdFromQuery(responseJSON.Query))
		} else {
			firstOccurrence := responseJSON.Occurrences[0]
			fmt.Println("Institution Code:", firstOccurrence.InstitutionCode)
			fmt.Println("Genus:", firstOccurrence.Genus)
			fmt.Println("CatalogNumber:", firstOccurrence.CatalogNumber)
			fmt.Println("Species:", firstOccurrence.Species)
			fmt.Println("StateProvince:", firstOccurrence.StateProvince)
			fmt.Println("Latitude:", firstOccurrence.Latitude)
			fmt.Println("DecimalLongitude:", firstOccurrence.Longitude)
			fmt.Println("Year:", firstOccurrence.Year)
			fmt.Println("Month:", firstOccurrence.MonthNumber)
			// fmt.Println("Day:", firstOccurrence.Day)
			fmt.Println("Collectors:", firstOccurrence.Collectors)
			fmt.Println("CollectorsNumber:", firstOccurrence.CollectorsNumber)

			day := getDayFromTimestamp(firstOccurrence.TimeStamp)
			fmt.Println("Day:", day)
		}

	}
}


func getDayFromTimestamp(timestamp int64) int {
	// Convert milliseconds to seconds for Go's time package
	eventTime := time.Unix(timestamp/1000, 0)
	// Extract the day of the month
	day := eventTime.Day()
	return day
}

func getIdFromQuery(query string) string {
	parts := strings.Split(query, ":")
	if len(parts) > 1 {
		result := strings.Trim(parts[1], "\" ")
		return result
	}
	return query
}

func createJson(result []byte) Response {
	var response Response
	// Unmarshal JSON data into the Repsonse struct
	err := json.Unmarshal(result, &response)
	if err != nil {
		fmt.Println("Error creating json:", err)
	}
	return response
}

func searchForCollections(collectionIds []string, urlStart string, urlEnd string) [][]byte {
	var allResponses [][]byte
	for _, collectionId := range collectionIds {
		fmt.Println(collectionId)
		url := urlStart + collectionId + urlEnd
		var response []byte = curl(url)
		allResponses = append(allResponses, response)
	}
	return allResponses
}

func curl(urlRequest string) []byte {
	resp, err := http.Get(urlRequest)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	return body
}
