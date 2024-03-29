package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "tenhoFibra/internal/geoHandler"
    "tenhoFibra/internal/prettify"
    "tenhoFibra/internal/protobuf"
)

type PostAddress struct {
    Message string 
}

func alive(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Yes, I'm Alive\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    log.Println("HTTP Post made - starting handlePostRquest")
    var addressData PostAddress 
	err := json.NewDecoder(r.Body).Decode(&addressData)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	address := addressData.Message
    coordinates := geoHandler.GetCoordinates(address)
    token := geoHandler.GetToken(coordinates.Y, coordinates.X)
    log.Println(token)
    objectsId := protobuf.DecodeProtobuf(token)
    log.Println(objectsId)
    data := geoHandler.GetNetworkInfo(objectsId)
    result := prettify.PrettifyData(data)
    log.Println("\n")
    log.Println(result)

    ids := geoHandler.GetProvidersIds()
    log.Println("\n")
    log.Println(ids)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": "Successful"})
}

func Server() {
    log.Println("Server Starting")
    http.HandleFunc("/alive", alive)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/tenhofibra", handlePostRequest)

    log.Println("Server Listening")
    http.ListenAndServe(":8090", nil)
}
