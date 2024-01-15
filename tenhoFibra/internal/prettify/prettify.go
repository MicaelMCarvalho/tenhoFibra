package prettify

import (
	//"encoding/base64"
	//"encoding/json"
	//"fmt"
	//"io"
	//"log"
	//"net/http"
	//"strconv"
	//"strings"
	"log"
	"tenhoFibra/internal/geoHandler"
)

type tenhoFibra struct {
}

func PrettifyData(input geoHandler.Info) tenhoFibra {
    var result tenhoFibra
    log.Println(input.MobileInfoData.RecordGroups[len(input.MobileInfoData.RecordGroups)-1])
    log.Println("\n")
    log.Println(input.TerrestrialInfoData.RecordGroups[len(input.TerrestrialInfoData.RecordGroups)-1])
    return result 
}

