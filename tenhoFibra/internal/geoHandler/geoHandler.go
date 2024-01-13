package geoHandler

import (
    "encoding/base64"
	"encoding/json"
    "fmt"
    "io"
	"log"
	"net/http"
    "strconv"
	"strings"
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

type coordinatesResponseBody struct {
	SpatialReference struct {
		Wkid       int `json:"wkid"`
		LatestWkid int `json:"latestWkid"`
	} `json:"spatialReference"`
	Candidates []Candidate `json:"candidates"`
}

type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Alias string `json:"alias"`
}

type Attribute struct {
	ObjectID   int `json:"objectid"`
	Operador   int `json:"operador"`
	VelDL7G    *int `json:"vel_dl_7g,omitempty"`
	VelDL6G    *int `json:"vel_dl_6g,omitempty"`
	VelDL5G    int `json:"vel_dl_5g"`
	VelDL4G    int `json:"vel_dl_4g"`
	VelDL3G    int `json:"vel_dl_3g"`
	VozMS3G    int `json:"voz_ms_3g"`
	VozMS2G    int `json:"voz_ms_2g"`
}

type NetworkInfoData struct {
	Field []Field `json:"fields"`
	RecordGroups []Attribute `json:"relatedRecordGroups"`
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

func getUrlNetInfo(objectIds []string) string {
    geocodeUrl := []string{
        "https://geo.anacom.pt/server/rest/services/publico/Coberturas_Disponiveis/MapServer/3/queryRelatedRecords?f=json&objectIds=",
        "&orderByFields=vel_dl_7g%20DESC%2Cvel_dl_6g%20DESC%2Cvel_dl_5g%20DESC%2Cvel_dl_4g%20DESC%2Cvel_dl_3g%20DESC%2Cvoz_ms_3g%20DESC%2Cvoz_ms_2g%20DESC&outFields=objectid%2Coperador%2Cvel_dl_7g%2Cvel_dl_6g%2Cvel_dl_5g%2Cvel_dl_4g%2Cvel_dl_3g%2Cvoz_ms_3g%2Cvoz_ms_2g&relationshipId=1&returnGeometry=true&definitionExpression=(vel_dl_7g%20is%20not%20null%20and%20vel_dl_7g%20%3C%3E%200)%20or%20(vel_dl_6g%20is%20not%20null%20and%20vel_dl_6g%20%3C%3E%200)%20or%20(vel_dl_5g%20is%20not%20null%20and%20vel_dl_5g%20%3C%3E%200)%20or%20(vel_dl_4g%20is%20not%20null%20and%20vel_dl_4g%20%3C%3E%200)%20or%20(vel_dl_3g%20is%20not%20null%20and%20vel_dl_3g%20%3C%3E%200)%20or%20(voz_ms_3g%20is%20not%20null%20and%20voz_ms_3g%20%3C%3E%200)%20or%20(voz_ms_2g%20is%20not%20null%20and%20voz_ms_2g%20%3C%3E%200)",
    }
    return geocodeUrl[0] + strings.Join(objectIds, ",") + geocodeUrl[1]
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
	var responseBody coordinatesResponseBody
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

func GetToken(y float64, x float64) string{
    url := getUrlToken(y, x)
    resp, err := http.Get(url)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return "" 
	}
    defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}
	base64Encoded := base64.StdEncoding.EncodeToString(body)

    return base64Encoded 
}

func GetNetworkInfo( ObjectsId []string ) NetworkInfoData {
    var result NetworkInfoData 
    url := getUrlNetInfo(ObjectsId)
    resp, err := http.Get(url)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return result
	}
    defer resp.Body.Close()

	var responseBody NetworkInfoData
    err1 := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err1 != nil {
		log.Fatal(err)
        return result
	}
    return responseBody
}
