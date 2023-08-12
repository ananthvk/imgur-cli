package internal

import (
	"encoding/json"
	"io"
)

type ImgData struct {
	Link       string `json:"link"`
	DeleteHash string `json:"deletehash"`
	Id         string `json:"id"`    // ID of the uploaded image
	Error      string `json:"error"` // Only available when Success is false
}

type Response struct {
	Data    ImgData `json:"data"`    // Data, depends on Success
	Success bool    `json:"success"` // Whether the request failed or not
	Status  int     `json:"status"`  // HTTP status
}

// This function decodes the response from imgur and populates the Response struct
func DecodeResponse(r io.Reader) (response Response, err error) {
	err = json.NewDecoder(r).Decode(&response)
	return
}
