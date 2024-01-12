package geoHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func getUrl(address string) string {
    geocodeUrl := []string{
        "https://geocode.arcgis.com/arcgis/rest/services/World/GeocodeServer/findAddressCandidates?f=json&outSR=%7B%22wkid%22:4326%7D&sourceCountry=PT&SingleLine=",
        "&maxSuggestions=1",
    }
    return geocodeUrl[0] + strings.Replace(address, " ", "%20", -1) + geocodeUrl[1]
}

func GetAddress(address string) string{
    log.Println("Get Address")
    url := getUrl(address)
    log.Println(url + "\n")

    resp, err := http.Get(url)
	if err != nil {
        //error handling
		return "Got an error doing http get to geocode.arcgis.com"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
        //Handle error
		return "ERROR"
	}

	response := map[string]interface{}{
		"status":     resp.Status,
		"statusCode": resp.StatusCode,
		"headers":    resp.Header,
		"body":       string(body),
	}

	// Convert the custom response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
        //Error Handling
		return "Error"
	}

	fmt.Println(string(jsonData))

    result := string("test")
    return result
}
