package geoHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
    "strconv"
)

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Candidate struct {
	Address    string   `json:"address"`
	Location   Location `json:"location"`
	Score      int      `json:"score"`
	Attributes struct{} `json:"attributes"`
	Extent     struct {
		Xmin float64 `json:"xmin"`
		Ymin float64 `json:"ymin"`
		Xmax float64 `json:"xmax"`
		Ymax float64 `json:"ymax"`
	} `json:"extent"`
}

type ResponseBody struct {
	SpatialReference struct {
		Wkid       int `json:"wkid"`
		LatestWkid int `json:"latestWkid"`
	} `json:"spatialReference"`
	Candidates []Candidate `json:"candidates"`
}


func getUrl(address string) string {
    geocodeUrl := []string{
        "https://geocode.arcgis.com/arcgis/rest/services/World/GeocodeServer/findAddressCandidates?f=json&outSR=%7B%22wkid%22:4326%7D&sourceCountry=PT&SingleLine=",
        "&maxSuggestions=1",
    }
    return geocodeUrl[0] + strings.Replace(address, " ", "%20", -1) + geocodeUrl[1]
}

func getUrlToken(y float64, x float64) string {
    geocodeUrl := []string{
        "https://geo.anacom.pt/server/rest/services/publico/Coberturas_Disponiveis/MapServer/3/query?f=pbf&geometry=%7B%22x%22%3A-7.810035006986%2C%22y%22%3A41.100895007188%7D&resultRecordCount=1&where=1%3D1&outFields=objectid&returnGeometry=false&spatialRel=esriSpatialRelIntersects&geometryType=esriGeometryPoint&inSR=4326",
        "%2C%22y%22%3A",
        "%7D&resultRecordCount=1&where=1%3D1&outFields=objectid&returnGeometry=false&spatialRel=esriSpatialRelIntersects&geometryType=esriGeometryPoint&inSR=4326",
    }
    return geocodeUrl[0] + strconv.FormatFloat(y, 'f', -1, 64) + geocodeUrl[1] + strconv.FormatFloat(x, 'f', -1, 64) + geocodeUrl[2]
}

func GetCoordinates(address string) Location {
    var location Location
    var locations []Location
    url := getUrl(address)

    resp, err := http.Get(url)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return location
	}

	var responseBody ResponseBody
    err1 := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err1 != nil {
		log.Fatal(err)
	}

    for _, candidate := range responseBody.Candidates {
        location := Location{
            X: candidate.Location.X,
            Y: candidate.Location.Y,
        }
        locations = append(locations, location)
    }
    return locations[0]
}


func GetToken(y float64, y float64) string{
    url := getUrlToken(y, x)

    resp, err := http.Get(url)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return "" 
	}
    log.Println(resp)

	// var responseBody ResponseBody
 //    err1 := json.NewDecoder(resp.Body).Decode(&responseBody)
	// if err1 != nil {
	// 	log.Fatal(err)
	// }
	//
    // for _, candidate := range responseBody.Candidates {
    //     location := Location{
    //         X: candidate.Location.X,
    //         Y: candidate.Location.Y,
    //     }
    //     locations = append(locations, location)
    // }
    return ""
}

