// data_service is used to make http requests to the dog.ceo API.
package data_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// getRandomImageAPIResponse is the response from the API.
// You can find the response structure in the API documentation.
// https://dog.ceo/dog-api/documentation/random
type getRandomImageAPIResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code,omitempty"`
}

// Creates a new http client with a timeout of 5 seconds.
// I would not create a new client for every request, but in this example it should be fine.
func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 5,
	}
}

// GetBreedImageURL returns the image URL as a string and an error if any.
func GetBreedImageURL(ctx context.Context, client *http.Client, breed, subBreed string) (string, int, error) {
	endpoint := createEndpoint(breed, subBreed)
	return getRandomImageURL(ctx, client, endpoint)
}

// GetImage returns the image as a byte array and an error if any.
// It downloads the image from the given URL.
func GetImage(ctx context.Context, client *http.Client, imageURL string) ([]byte, int, error) {
	return processHttpGet(ctx, client, imageURL)
}

// createEndpoint returns the endpoint URL for the given breed and sub-breed.
// If the sub-breed is empty, it returns the endpoint for the breed.
// If the sub-breed is not empty, it returns the endpoint for breed and the sub-breed.
// Example:
// breed: "husky"
// subBreed: ""
// endpoint: "https://dog.ceo/api/breed/husky/images/random"
func createEndpoint(breed, subBreed string) string {
	if len(subBreed) > 0 {
		return fmt.Sprintf("https://dog.ceo/api/breed/%s/%s/images/random", breed, subBreed)
	}
	return fmt.Sprintf("https://dog.ceo/api/breed/%s/images/random", breed)
}

// getRandomImageURL returns the image URL as a string, status code as an integer and an error if any.
// It uses the given endpoint to get the image URL.
func getRandomImageURL(ctx context.Context, client *http.Client, endpoint string) (string, int, error) {
	resp, statusCode, err := processHttpGet(ctx, client, endpoint)
	if err != nil {
		return "", statusCode, err
	}

	if statusCode != http.StatusOK {
		return "", statusCode, nil
	}

	apiResp := &getRandomImageAPIResponse{}
	if err := json.Unmarshal(resp, apiResp); err != nil {
		return "", statusCode, err
	}

	return string(apiResp.Message), statusCode, nil
}

// processHttpGet returns the response as a byte array, status code as an integer and an error if any.
// It uses the given endpoint to get the response.
func processHttpGet(ctx context.Context, client *http.Client, endpoint string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return body, resp.StatusCode, nil
}
