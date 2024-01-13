package protobuf

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Response struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

var url string = "http://localhost:3000/decode"

func DecodeProtobuf(token string) []string {
    var result []string

	data := map[string]string{
		"token": token,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return result
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return result 
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return result 
	}

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result 
	}
    log.Println(body)

    var response Response
	err = json.Unmarshal(body.Bytes(), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result 
	}

	if response.Status != 200 {
		fmt.Println("Unexpected status code:", response.Status)
		return result 
	}

	return response.Data
}
