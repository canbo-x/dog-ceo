package breed_image_service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/canbo-x/dog-ceo/data_service"
)

// statusNotFoundErrorText is the error text for the status code 404
var errStatusNotFound = errors.New("image is not found on the server! Please check the url or search again")

// GetImage fetch the image from the given url and returns the image bytes and an error if any
func GetImage(ctx context.Context, client *http.Client, imageURL string) ([]byte, error) {
	image, statusCode, err := data_service.GetImage(ctx, client, imageURL)
	if err != nil {
		return nil, err
	}
	if statusCode == http.StatusNotFound {
		return nil, errStatusNotFound
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with : %d", statusCode)
	}
	return image, nil
}

// GetURL returns the image URL as a string and an error if any.
// It throws an error if the status code is not 200.
func GetURL(ctx context.Context, client *http.Client, breed string, subBreed string) (string, error) {
	imageURL, statusCode, err := data_service.GetBreedImageURL(ctx, client, breed, subBreed)
	if err != nil {
		return imageURL, err
	}
	if statusCode == http.StatusNotFound {
		return imageURL, errStatusNotFound
	}
	if statusCode != http.StatusOK {
		return imageURL, fmt.Errorf("server responded with : %d", statusCode)
	}
	return imageURL, err
}
