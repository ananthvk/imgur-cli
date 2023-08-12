package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var logger = log.New(os.Stdout, "", 0)
var wg sync.WaitGroup

// Upload - Uploads the list of files to imgur.
// files - Slice of strings, which represents paths of the files to upload

func uploadFile(file string, client *http.Client) {
	defer wg.Done()
	// Create a buffer for storing the multipart data
	buffer := new(bytes.Buffer)
	// Create a multipart writer
	writer := multipart.NewWriter(buffer)
	imgWriter, err := writer.CreateFormFile("image", filepath.Base(file))
	if err != nil {
		log.Println("internal error: creating writer", file)
		log.Fatal(err)
	}
	// Read the file from disk
	fil, err := os.Open(file)
	if err != nil {
		log.Println("error: could not open ", file)
		log.Fatal(err)
	}
	defer fil.Close()
	// Copy the file to the buffer through imgWriter
	_, err = io.Copy(imgWriter, fil)
	if err != nil {
		log.Println("internal error: while copying", file)
		log.Fatal(err)
	}
	writer.Close()
	// Create the headers, and fields
	req, err := http.NewRequest(http.MethodPost, API_URL, buffer)
	if err != nil {
		log.Println("internal error: while creating http request")
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Client-Id "+os.Getenv(CLIENT_ID_ENV_NAME))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(req)
	if err != nil {
		log.Println("error: could not upload file to imgur")
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Read the response and decode the json into the struct

	resp, err := DecodeResponse(response.Body)
	if err != nil {
		log.Println("error: could not read response")
		log.Fatal(err)
	}
	// Handle HTTP errors
	switch resp.Status {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNonAuthoritativeInfo:
		logger.Println("Success: uploaded image to imgur!")
		logger.Println("Image URL:", resp.Data.Link)
		logger.Println("Delete hash:", resp.Data.DeleteHash)
		logger.Println("Please keep the above delete hash safe as it is required to remove the image from imgur.")
	case http.StatusUnauthorized:
		log.Println("Unauthorized")
		log.Fatal("Please check if you have set a valid", CLIENT_ID_ENV_NAME)
	default:
		log.Println("Error:", resp.Status)
		log.Fatal(resp.Data.ErrorString)
	}
}
func Upload(files []string) {
	if len(strings.TrimSpace(os.Getenv(CLIENT_ID_ENV_NAME))) == 0 {
		fmt.Fprintf(os.Stderr, HELP_ENV_MESSAGE)
		os.Exit(2)
	}
	client := &http.Client{Timeout: TIMEOUT}
	for _, file := range files {
		wg.Add(1)
		go uploadFile(file, client)
	}
	wg.Wait()
}
