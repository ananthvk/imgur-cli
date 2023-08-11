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
	"time"
)

const API_URL string = "https://api.imgur.com/3/image"
const CLIENT_ID_ENV_NAME = "IMGUR_CLIENT_ID"
const TIMEOUT = 15 * time.Second
const HELP_ENV_MESSAGE = "This program requires a client id to upload images to imgur.\n\n" +
	CLIENT_ID_ENV_NAME + ` environment variable not set.
	Get a client id from https://api.imgur.com/oauth2/addclient

	If you are on windows,
	set ` + CLIENT_ID_ENV_NAME + "=[CLIENT_ID]" + `

	If you are on linux,
	export ` + CLIENT_ID_ENV_NAME + "=[CLIENT_ID]\n"

// Upload - Uploads the list of files to imgur.
// files - Slice of strings, which represents paths of the files to upload

func uploadFile(file string) {
	// Create a buffer for storing the multipart data
	buffer := new(bytes.Buffer)
	// Create a multipart writer
	writer := multipart.NewWriter(buffer)
	imgWriter, err := writer.CreateFormFile("image", filepath.Base(file))
	if err != nil {
		log.Print("internal error: creating writer", file)
		log.Fatal(err)
	}
	// Read the file from disk
	fil, err := os.Open(file)
	if err != nil {
		log.Print("error: could not open ", file)
		log.Fatal(err)
	}
	defer fil.Close()
	// Copy the file to the buffer through imgWriter
	_, err = io.Copy(imgWriter, fil)
	if err != nil {
		log.Print("internal error: while copying", file)
		log.Fatal(err)
	}
	writer.Close()
	// Create the headers, and fields
	client := &http.Client{Timeout: TIMEOUT}
	req, err := http.NewRequest(http.MethodPost, API_URL, buffer)
	if err != nil {
		log.Print("internal error: while creating http request")
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Client-Id "+os.Getenv(CLIENT_ID_ENV_NAME))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(req)
	if err != nil {
		log.Print("error: could not upload file to imgur")
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print("error: could not read response")
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(body))
}
func Upload(files []string) {
	if len(strings.TrimSpace(os.Getenv(CLIENT_ID_ENV_NAME))) == 0 {
		fmt.Fprintf(os.Stderr, HELP_ENV_MESSAGE)
		os.Exit(2)
	}
	for _, file := range files {
		uploadFile(file)
	}
}
