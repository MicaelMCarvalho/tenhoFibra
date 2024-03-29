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

type AttributeMobile struct {
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

type RelatedRecordsMobile struct {
	Attributes AttributeMobile `json:"attributes"`
}

type RecordGroupsMobile struct {
	ObjectID   int `json:"objectId"`
    RelatedRecords []RelatedRecordsMobile `json:"relatedRecords"`
}

type MobileInfoData struct {
    Field []Field `json:"fields"`
    RecordGroups []RecordGroupsMobile `json:"relatedRecordGroups"`
}


type AttributeTerrestrial struct {
	ObjectID   int `json:"objectid"`
	Operador   int `json:"operador"`
    TechnA     int `json:"tecnologia_a"`
	VelDwA     int `json:"vel_max_dl_a"`
	VelUpA     int `json:"vel_max_ul_a"`
    TechnB     int `json:"tecnologia_b"`
	VelDwB     int `json:"vel_max_dl_b"`
	VelUp     int `json:"vel_max_ul_b"`
}

type RelatedRecordsTerrestrial struct {
	Attributes AttributeTerrestrial `json:"attributes"`
}

type RecordGroupsTerrestrial struct {
	ObjectID   int `json:"objectId"`
    RelatedRecords []RelatedRecordsTerrestrial `json:"relatedRecords"`
}

type TerrestrialInfoData struct {
	Field []Field `json:"fields"`
	RecordGroups []RecordGroupsTerrestrial `json:"relatedRecordGroups"`
}

type Info struct {
    MobileInfoData MobileInfoData
    TerrestrialInfoData TerrestrialInfoData
}

type Provider struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	ParentLayerID         int      `json:"parentLayerId"`
	DefaultVisibility     bool     `json:"defaultVisibility"`
	SubLayerIds           []int    `json:"subLayerIds"`
	MinScale              int      `json:"minScale"`
	MaxScale              int      `json:"maxScale"`
	Type                  string   `json:"type"`
	SupportsDynamicLegends bool     `json:"supportsDynamicLegends"`
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

func getUrlMobileInfo(objectIds []string) string {
    geocodeUrl := []string{
        "https://geo.anacom.pt/server/rest/services/publico/Coberturas_Disponiveis/MapServer/3/queryRelatedRecords?f=json&objectIds=",
        "&orderByFields=vel_dl_7g%20DESC%2Cvel_dl_6g%20DESC%2Cvel_dl_5g%20DESC%2Cvel_dl_4g%20DESC%2Cvel_dl_3g%20DESC%2Cvoz_ms_3g%20DESC%2Cvoz_ms_2g%20DESC&outFields=objectid%2Coperador%2Cvel_dl_7g%2Cvel_dl_6g%2Cvel_dl_5g%2Cvel_dl_4g%2Cvel_dl_3g%2Cvoz_ms_3g%2Cvoz_ms_2g&relationshipId=1&returnGeometry=true&definitionExpression=(vel_dl_7g%20is%20not%20null%20and%20vel_dl_7g%20%3C%3E%200)%20or%20(vel_dl_6g%20is%20not%20null%20and%20vel_dl_6g%20%3C%3E%200)%20or%20(vel_dl_5g%20is%20not%20null%20and%20vel_dl_5g%20%3C%3E%200)%20or%20(vel_dl_4g%20is%20not%20null%20and%20vel_dl_4g%20%3C%3E%200)%20or%20(vel_dl_3g%20is%20not%20null%20and%20vel_dl_3g%20%3C%3E%200)%20or%20(voz_ms_3g%20is%20not%20null%20and%20voz_ms_3g%20%3C%3E%200)%20or%20(voz_ms_2g%20is%20not%20null%20and%20voz_ms_2g%20%3C%3E%200)",
    }
    return geocodeUrl[0] + strings.Join(objectIds, ",") + geocodeUrl[1]
}

func getUrlTerrestrialInfo(objectIds []string) string {
    geocodeUrl := []string{
        "https://geo.anacom.pt/server/rest/services/publico/Coberturas_Disponiveis/MapServer/0/queryRelatedRecords?f=json&objectIds=",
        "&orderByFields=vel_max_dl_a%20DESC%2Cvel_max_dl_b%20DESC%2Cvel_max_ul_a%20DESC%2Cvel_max_ul_b%20DESC&outFields=objectid%2Coperador%2Ctecnologia_a%2Cvel_max_dl_a%2Cvel_max_ul_a%2Ctecnologia_b%2Cvel_max_dl_b%2Cvel_max_ul_b&relationshipId=0&returnGeometry=true&definitionExpression=((vel_max_dl_a%20is%20not%20null%20and%20vel_max_dl_a%20%3C%3E%200)%20and%20(vel_max_ul_a%20is%20not%20null%20and%20vel_max_ul_a%20%3C%3E%200))%20or%20((vel_max_dl_b%20is%20not%20null%20and%20vel_max_dl_b%20%3C%3E%200)%20and%20(vel_max_ul_b%20is%20not%20null%20and%20vel_max_ul_b%20%3C%3E%200))",
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

func GetNetworkInfo(ObjectsId []string) Info {
    var result Info
    urlMobile := getUrlMobileInfo(ObjectsId)
    respMobile, err := http.Get(urlMobile)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return result
	}
    defer respMobile.Body.Close()

	var dataMobile MobileInfoData 
    err1 := json.NewDecoder(respMobile.Body).Decode(&dataMobile)
	if err1 != nil {
		log.Fatal(err1)
        return result
	}

    urlTerrestrial := getUrlTerrestrialInfo(ObjectsId)
    respTerrestrial, err := http.Get(urlTerrestrial)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return result
	}
    defer respTerrestrial.Body.Close()

	var dataTerrestrial TerrestrialInfoData
    err2 := json.NewDecoder(respTerrestrial.Body).Decode(&dataTerrestrial)
	if err2 != nil {
		log.Fatal(err2)
        return result
	}

    result.MobileInfoData = dataMobile
    result.TerrestrialInfoData = dataTerrestrial
    
    return result 
}


func GetProvidersIds() string {
    var result string 
    url := "https://geo.anacom.pt/server/rest/services/publico/EstatisticasMercado_Pub/MapServer?f=json"
    resp, err := http.Get(url)
	if err != nil {
        log.Println("Got an error doing http get to geocode.arcgis.com")
        return result
	}
    defer resp.Body.Close()
    
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return result
	}

	var info map[string]json.RawMessage

	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return result
	}

    log.Println("Reading json")
	if value, ok := info["layers"]; ok {
		var provider []Provider
		err := json.Unmarshal(value, &provider)
		if err != nil {
			fmt.Println("Error decoding field:", err)
		    return result
		}
		fmt.Println("Value of fieldName:", provider)
	} else {
		fmt.Println("Field fieldName not found in JSON")
	}

    return result 
}
