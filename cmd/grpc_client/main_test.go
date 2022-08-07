package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/canbo-x/dog-ceo/proto/breed_image"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type mockServer struct {
	breed_image.UnimplementedBreedImageServiceServer
}

func (*mockServer) Search(ctx context.Context, req *breed_image.BreedImageSearchRequest) (*breed_image.BreedImageSearchResponse, error) {
	return &breed_image.BreedImageSearchResponse{ImageURL: "test_url", Image: []byte("test")}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	breed_image.RegisterBreedImageServiceServer(server, &mockServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func getMockCoon() (context.Context, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return ctx, conn
}

func getRealCoon() (context.Context, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, getAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return ctx, conn
}

func TestClientWithMockServer(t *testing.T) {
	tests := map[string]struct {
		Breed    string
		SubBreed string
		Valid    bool
	}{
		"valid breed": {
			Breed:    "husky",
			SubBreed: "",
			Valid:    true,
		},
		"valid subbreed": {
			Breed:    "australian",
			SubBreed: "shepherd",
			Valid:    true,
		},
	}

	ctx, conn := getMockCoon()
	defer conn.Close()

	client := breed_image.NewBreedImageServiceClient(conn)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := client.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: test.Breed, SubBreed: test.SubBreed})
			if (err == nil) != test.Valid {
				t.Fatalf("error: %v", err)
			}
		})
	}
}

func TestClientWithRealServer(t *testing.T) {
	tests := map[string]struct {
		Breed           string
		SubBreed        string
		ShouldHaveURL   bool
		ShouldHaveImage bool
		Valid           bool
	}{
		"valid breed": {
			Breed:           "husky",
			SubBreed:        "",
			ShouldHaveURL:   true,
			ShouldHaveImage: true,
			Valid:           true,
		},
		"valid subbreed": {
			Breed:           "australian",
			SubBreed:        "shepherd",
			ShouldHaveURL:   true,
			ShouldHaveImage: true,
			Valid:           true,
		},
		"invalid breed": {
			Breed:           "invalid",
			SubBreed:        "",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"invalid subbreed": {
			Breed:           "australian",
			SubBreed:        "invalid",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"invalid breed and subbreed": {
			Breed:           "invalid",
			SubBreed:        "invalid",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"regexp breed": {
			Breed:           "BROKEN_REGEX",
			SubBreed:        "",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"regexp subbreed": {
			Breed:           "australian",
			SubBreed:        "BROKEN_REGEX",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"whitespace breed": {
			Breed:           " ",
			SubBreed:        "",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
		"whitespace subbreed": {
			Breed:           "australian",
			SubBreed:        " ",
			ShouldHaveURL:   false,
			ShouldHaveImage: false,
			Valid:           false,
		},
	}

	ctx, conn := getRealCoon()
	defer conn.Close()

	client := breed_image.NewBreedImageServiceClient(conn)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			resp, err := client.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: test.Breed, SubBreed: test.SubBreed})
			if (err == nil) != test.Valid {
				t.Fatalf("error: %v", err)
			}

			if test.ShouldHaveImage && len(resp.Image) == 0 {
				t.Fatalf("image is empty")
			}

			if test.ShouldHaveURL && len(resp.ImageURL) == 0 {
				t.Fatalf("url is empty")
			}

		})
	}

}

func TestGetAddr(t *testing.T) {
	addr := getAddr()
	if addr == "" {
		t.Fatalf("addr is empty")
	}

	t.Setenv("CLIENT_GRPC_ADDR", "localhost:12345")
	addr = getAddr()
	if addr != "localhost:12345" {
		t.Fatalf("addr is read correctly after set")
	}
}

func TestHandleFileName(t *testing.T) {
	const defaultURL = "https://images.dog.ceo/breeds/husky/default.jpg"
	const defaultExpectedFileName = "default.jpg"

	tests := map[string]struct {
		Name         string
		Url          string
		ExpectedName string
		Valid        bool
	}{
		"valid name": {
			Name:         "lovelyHusky",
			Url:          defaultURL,
			ExpectedName: "lovelyHusky.jpg",
			Valid:        true,
		},
		"valid name with number": {
			Name:         "lovelyHusky1",
			Url:          defaultURL,
			ExpectedName: "lovelyHusky1.jpg",
			Valid:        true,
		},
		"valid name with underscore": {
			Name:         "lovely_husky",
			Url:          defaultURL,
			ExpectedName: "lovely_husky.jpg",
			Valid:        true,
		},
		"valid name with dash": {
			Name:         "lovely-Husky",
			Url:          defaultURL,
			ExpectedName: "lovely-Husky.jpg",
			Valid:        true,
		},
		"valid name with number and underscore": {
			Name:         "lovely_husky1",
			Url:          defaultURL,
			ExpectedName: "lovely_husky1.jpg",
			Valid:        true,
		},
		"invalid name with dot": {
			Name:         "lovely.Husky",
			Url:          defaultURL,
			ExpectedName: "",
			Valid:        false,
		},
		"invalid name with space": {
			Name:         "lovely Husky",
			Url:          defaultURL,
			ExpectedName: "",
			Valid:        false,
		},
		"invalid name with special characters": {
			Name:         "lovely%Husky",
			Url:          defaultURL,
			ExpectedName: "",
			Valid:        false,
		},
		"empty name": {
			Name:         "",
			Url:          defaultURL,
			ExpectedName: defaultExpectedFileName,
			Valid:        true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := handleFileName(test.Name, test.Url)
			if (err == nil) != test.Valid || got != test.ExpectedName {
				t.Fatalf("want %s; got %s", test.ExpectedName, got)
			}
		})
	}

}
