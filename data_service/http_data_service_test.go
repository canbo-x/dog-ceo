package data_service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient()
	if reflect.TypeOf(client) != reflect.TypeOf(&http.Client{}) {
		t.Errorf("method : NewHttpClient does not return http.Client")
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 7)
	}))
	defer ts.Close()
	client := NewHttpClient()
	_, err := client.Get(ts.URL)
	if err != nil && !os.IsTimeout(err) {
		t.Errorf("timeout error is not returned")
	}
}

func TestGetBreedImageURL(t *testing.T) {
	tests := map[string]struct {
		Breed              string
		SubBreed           string
		ExpectedURL        string
		ExpectedStatusCode int
		Valid              bool
	}{
		"valid breed and subbreed": {
			Breed:              "australian",
			SubBreed:           "shepherd",
			ExpectedURL:        "https://images.dog.ceo/breeds/australian-shepherd/",
			ExpectedStatusCode: http.StatusOK,
			Valid:              true,
		},
		"valid breed and empty subbreed - dog does not have subbreed": {
			Breed:              "husky",
			SubBreed:           "",
			ExpectedURL:        "https://images.dog.ceo/breeds/husky/",
			ExpectedStatusCode: http.StatusOK,
			Valid:              true,
		},
		"valid breed and empty subbreed - dog have subbreed ": {
			Breed:              "australian",
			SubBreed:           "",
			ExpectedURL:        "https://images.dog.ceo/breeds/australian-shepherd/",
			ExpectedStatusCode: http.StatusOK,
			Valid:              true,
		},
		"invalid breed": {
			Breed:              "INVALID",
			SubBreed:           "",
			ExpectedURL:        "",
			ExpectedStatusCode: http.StatusNotFound,
			Valid:              true,
		},
		"invalid breed and subbreed": {
			Breed:              "INVALID",
			SubBreed:           "INVALID",
			ExpectedURL:        "",
			ExpectedStatusCode: http.StatusNotFound,
			Valid:              true,
		},
		"valid breed and invalid subbreed": {
			Breed:              "australian",
			SubBreed:           "INVALID",
			ExpectedURL:        "",
			ExpectedStatusCode: http.StatusNotFound,
			Valid:              true,
		},
		"empty breed and subbreed": {
			Breed:              "",
			SubBreed:           "",
			ExpectedURL:        "",
			ExpectedStatusCode: http.StatusNotFound,
			Valid:              true,
		},
		"broken url": {
			Breed:              "%INVALID._=",
			SubBreed:           "",
			ExpectedURL:        "",
			ExpectedStatusCode: http.StatusInternalServerError,
			Valid:              false,
		},
	}

	client := NewHttpClient()

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			url, statusCode, err := GetBreedImageURL(context.Background(), client, test.Breed, test.SubBreed)
			if (err == nil) != test.Valid {
				t.Fatalf("want err == nil => %t; got err %v", test.Valid, err)
			}
			if statusCode != test.ExpectedStatusCode {
				t.Fatalf("status code is not correct got %d want %d", statusCode, test.ExpectedStatusCode)
			}
			if !strings.Contains(url, test.ExpectedURL) {
				t.Fatalf("url is not correct got: %v expected to contain: %v", url, test.ExpectedURL)
			}
		})
	}
}

func TestGetImage(t *testing.T) {
	tests := map[string]struct {
		URL                string
		ExpectedStatusCode int
		ShouldHaveImage    bool
		Valid              bool
	}{
		"valid url": {
			URL:                "https://images.dog.ceo/breeds/husky/n02110185_12678.jpg",
			ExpectedStatusCode: http.StatusOK,
			ShouldHaveImage:    true,
			Valid:              true,
		},
		"invalid url": {
			URL:                "https://images.dog.ceo/breeds/INVALID/no.jpg",
			ExpectedStatusCode: http.StatusNotFound,
			ShouldHaveImage:    false,
			Valid:              true,
		},
	}

	client := NewHttpClient()
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			image, statusCode, err := GetImage(context.Background(), client, test.URL)
			if (err == nil) != test.Valid {
				t.Fatalf("want err == nil => %t; got err %v", test.Valid, err)
			}
			if statusCode != test.ExpectedStatusCode {
				t.Fatalf("status code is not correct got %d want %d", statusCode, test.ExpectedStatusCode)
			}
			if test.ShouldHaveImage && image == nil {
				t.Fatalf("image is not correct got: %v expected to contain: %v", image, test.ShouldHaveImage)
			}
		})
	}
}

func TestCreateEndpoint(t *testing.T) {
	tests := map[string]struct {
		Breed       string
		SubBreed    string
		ExpectedURL string
	}{
		"only breed": {
			Breed:       "australian",
			SubBreed:    "",
			ExpectedURL: "https://dog.ceo/api/breed/australian/images/random",
		},
		"empty breed and subbreed": {
			Breed:       "",
			SubBreed:    "",
			ExpectedURL: "https://dog.ceo/api/breed//images/random",
		},
		"both breed and subbreed": {
			Breed:       "australian",
			SubBreed:    "shepherd",
			ExpectedURL: "https://dog.ceo/api/breed/australian/shepherd/images/random",
		},
		"only subbreed": {
			Breed:       "",
			SubBreed:    "shepherd",
			ExpectedURL: "https://dog.ceo/api/breed//shepherd/images/random",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			url := createEndpoint(test.Breed, test.SubBreed)
			if url != test.ExpectedURL {
				t.Fatalf("want url %v; got %v", test.ExpectedURL, url)
			}
		})
	}

}
