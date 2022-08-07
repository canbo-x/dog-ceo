package main

import (
	"context"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/canbo-x/dog-ceo/dummy_rate_limiter"
	"github.com/canbo-x/dog-ceo/proto/breed_image"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(shouldLimit bool) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpcServer()
	if shouldLimit {
		server = grpcServerWithRateLimit()
	}

	breed_image.RegisterBreedImageServiceServer(server, &breedImageServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func grpcServer() *grpc.Server {
	return grpc.NewServer()
}

func grpcServerWithRateLimit() *grpc.Server {
	dummyRL := dummy_rate_limiter.NewLimitCounter()
	dummyRL.StartLimiter()
	return grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			ratelimit.UnaryServerInterceptor(dummyRL),
		),
	)
}

func getCoon(shouldLimit bool) (context.Context, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer(shouldLimit)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return ctx, conn
}

func getClient(conn *grpc.ClientConn) breed_image.BreedImageServiceClient {
	return breed_image.NewBreedImageServiceClient(conn)
}
func TestServerAndSearch(t *testing.T) {
	tests := map[string]struct {
		Breed    string
		SubBreed string
		URL      string
		Valid    bool
	}{
		"valid breed and subbreed": {
			Breed:    "australian",
			SubBreed: "shepherd",
			URL:      "https://images.dog.ceo/breeds/australian-shepherd/",
			Valid:    true,
		},
		"valid breed and empty subbreed - dog does not have subbreed": {
			Breed:    "husky",
			SubBreed: "",
			URL:      "https://images.dog.ceo/breeds/husky/",
			Valid:    true,
		},
		"valid breed and empty subbreed - dog have subbreed ": {
			Breed:    "australian",
			SubBreed: "",
			URL:      "https://images.dog.ceo/breeds/australian-shepherd/",
			Valid:    true,
		},
		"invalid breed": {
			Breed:    "INVALID",
			SubBreed: "",
			URL:      "",
			Valid:    false,
		},
		"invalid subbreed": {
			Breed:    "husky",
			SubBreed: "INVALID",
			URL:      "",
			Valid:    false,
		},
		"empty breed": {
			Breed:    "",
			SubBreed: "",
			URL:      "",
			Valid:    false,
		},
		"breed regex not match": {
			Breed:    "INVALID_REGEX",
			SubBreed: "",
			URL:      "",
			Valid:    false,
		},
		"subbreed regex not match": {
			Breed:    "australian",
			SubBreed: "INVALID_REGEX",
			URL:      "",
			Valid:    false,
		},
		"whitespace breed": {
			Breed:    " ",
			SubBreed: "",
			URL:      "",
			Valid:    false,
		},
		"whitespace subbreed": {
			Breed:    "australian",
			SubBreed: " ",
			URL:      "",
			Valid:    false,
		},
	}

	ctx, conn := getCoon(false)
	defer conn.Close()
	client := getClient(conn)

	for name, test := range tests {
		test := test
		// https://github.com/golang/go/issues/17791#issuecomment-259976524
		t.Run("group", func(t *testing.T) {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				resp, err := client.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: test.Breed, SubBreed: test.SubBreed})
				if (err == nil) != test.Valid {
					t.Fatalf("want err == nil => %t; got err %v", test.Valid, err)
				}

				if !resp.ProtoReflect().IsValid() {
					if err == nil {
						t.Fatalf("response is invalid and there is no error")
					}
					// if we are already here it means:
					// error is not nil which means response must be nil
					// to avoid invalid memory address or nil pointer dereference error
					// we return
					return
				}

				if !strings.Contains(resp.ImageURL, test.URL) && test.Valid {
					t.Fatalf("want url contains %v; got %v", test.URL, resp.ImageURL)
				}

				if (resp.Image == nil) == test.Valid {
					t.Fatalf("image is nil")
				}
			})
		})
	}
}

func TestRateLimiting(t *testing.T) {
	ctx, conn := getCoon(true)
	defer conn.Close()
	client := getClient(conn)

	for i := 0; i < 10; i++ {
		_, err := client.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: "australian", SubBreed: "shepherd"})
		if err != nil && status.Code(err) != codes.ResourceExhausted {
			t.Fatalf("error is not ResourceExhausted %v", err)
		}
	}

	time.Sleep(3 * time.Second)
	_, err := client.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: "australian", SubBreed: "shepherd"})
	if err != nil {
		t.Fatalf("error is not nil %v", err)
	}

}
