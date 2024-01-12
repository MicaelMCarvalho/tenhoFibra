package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "tenhoFibra/internal/geoHandler"
)

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

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
    result := geoHandler.GetCoordinates(address)
    log.Println(result.Y)
    log.Println(result.x)
    result_token := geoHandler.GetToken(result.Y, result.X)
    log.Println(result_token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": "Successful"})
}

func Server() {
    log.Println("Server Starting")
    http.HandleFunc("/alive", alive)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/tenhofibra", handlePostRequest)

    http.ListenAndServe(":8090", nil)
}
