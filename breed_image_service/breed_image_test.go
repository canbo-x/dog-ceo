package breed_image_service

import (
	"context"
	"net/http"
	"testing"
)

func TestGetImage(t *testing.T) {
	image, err := GetImage(context.Background(), &http.Client{}, "https://images.dog.ceo/breeds/husky/n02110185_5030.jpg")
	if err != nil {
		t.Error(err)
	}
	if len(image) == 0 {
		t.Error("image is empty")
	}
}

func TestGetImageNotFound(t *testing.T) {
	_, err := GetImage(context.Background(), &http.Client{}, "broken_link")
	if err == nil {
		t.Error("image is not found")
	}
}

func TestGetURL(t *testing.T) {
	url, err := GetURL(context.Background(), &http.Client{}, "husky", "")
	if err != nil {
		t.Error(err)
	}
	if url == "" {
		t.Error("url is empty")
	}
}

func TestGetURLNotFound(t *testing.T) {
	_, err := GetURL(context.Background(), &http.Client{}, "husky", "not-found")
	if err == nil {
		t.Error("url is not found")
	}
}

func TestGetURLInvalid(t *testing.T) {
	_, err := GetURL(context.Background(), &http.Client{}, "", "")
	if err == nil {
		t.Error("url is invalid")
	}
}
