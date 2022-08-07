package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/canbo-x/dog-ceo/breed_image_service"
	"github.com/canbo-x/dog-ceo/data_service"
	"github.com/canbo-x/dog-ceo/dummy_rate_limiter"
	"github.com/canbo-x/dog-ceo/proto/breed_image"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var isValidString = regexp.MustCompile(`^[A-Za-z]+$`).MatchString

// Implement the breed image server.
type breedImageServer struct {
	breed_image.UnimplementedBreedImageServiceServer
}

// Search checks for the image of the given breed and sub-breed.
func (bis *breedImageServer) Search(ctx context.Context, bi *breed_image.BreedImageSearchRequest) (*breed_image.BreedImageSearchResponse, error) {
	log.Printf("Received a request to search. Breed : %v Sub Breed : %v\n", bi.Breed, bi.SubBreed)

	if !isValidString(bi.Breed) {
		log.Println("Invalid breed name. Request is rejected.")
		return nil, fmt.Errorf("invalid breed name it can only contains english latin letters : %v", bi.Breed)
	}

	if bi.SubBreed != "" && !isValidString(bi.SubBreed) {
		log.Println("Invalid sub-breed name. Request is rejected.")
		return nil, fmt.Errorf("invalid sub-breed name it can only contains english latin letters : %v", bi.SubBreed)
	}

	imageURL, err := breed_image_service.GetURL(ctx, data_service.NewHttpClient(), bi.Breed, bi.SubBreed)
	if err != nil {
		log.Printf("Error while getting image url : %v\n", err)
		return nil, err
	}

	image, err := breed_image_service.GetImage(ctx, data_service.NewHttpClient(), imageURL)
	if err != nil {
		log.Printf("Error while getting image : %v\n", err)
		return nil, fmt.Errorf("failed to get image : %v", err)
	}

	log.Printf("Image is fetched and served to the client. Image URL : %v\n", imageURL)
	return &breed_image.BreedImageSearchResponse{ImageURL: imageURL, Image: image}, nil
}

func main() {

	// This port is used to serve the breed image service.
	port := flag.Int("port", 22626, "The gRPC-server port.")

	// This log level is used to set the log level.
	logLevel := flag.String("log-level", "info", "The log level of the gRPC-server.")

	// Parse the command line flags
	flag.Parse()

	logrusLogger := logrus.New()
	if err := checkAndSetLogLevel(*logLevel); err != nil {
		logrusLogger.Fatalf("Failed to set log level : %v", err)
	}

	dummyRL := dummy_rate_limiter.NewLimitCounter()
	dummyRL.StartLimiter()

	// Listen on the port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrusLogger.Fatalf("failed to listen: %v", err)
	}

	logrusEntry := logrus.NewEntry(logrusLogger)
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			ratelimit.UnaryServerInterceptor(dummyRL),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			ratelimit.StreamServerInterceptor(dummyRL),
		),
	)

	// Register the breed image server
	breed_image.RegisterBreedImageServiceServer(server, &breedImageServer{})
	logrusLogger.Infof("gRPC server is listening on port %d", *port)

	errChan := make(chan error)

	// use buffered channel
	// https://pkg.go.dev/os/signal#Notify and https://gobyexample.com/signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	defer func() {
		server.GracefulStop()
	}()

	select {
	case err := <-errChan:
		logrusLogger.Fatalf("Fatal error: %v\n", err)
	case <-stopChan:
		logrusLogger.Info("Stopping the server...")
	}
}

// checkAndSetLogLevel checks the log level and sets it.
// If the log level is invalid, it returns an error.
func checkAndSetLogLevel(logLevel string) error {
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		return fmt.Errorf("invalid log level : %v", logLevel)
	}
	return nil
}
