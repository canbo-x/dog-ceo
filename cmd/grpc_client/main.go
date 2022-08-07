package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/canbo-x/dog-ceo/proto/breed_image"
	"github.com/canbo-x/dog-ceo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	help := flag.Bool("help", false, "flag to show help")
	flag.Parse()

	if *help {
		helpCommand()
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(getAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Set up the client
	c := breed_image.NewBreedImageServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(os.Args) < 2 {
		log.Println("expected a command please run `<executable> help` for more information")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "search":
		searchCommand(ctx, c, os.Args[2:])
	default:
		log.Println("expected a valid command please run `<executable> help` for more information")
		os.Exit(1)
	}
}

// searchCommand searches for the breed image.
// Breed is required.
// Sub-breed, save, path and file-name are optional.
// If the sub-breed is provided, it searches for the sub-breed image.
// If the save flag is provided, it saves the image to the path.
// If the path is not provided, it saves the image to the default directory [images/].
// If the file name is not provided, it gets the file name from the image URL.
// If save flag is provided, it prints the full path of the image.
// If save flag is not provided, it prints the image URL only.
func searchCommand(ctx context.Context, c breed_image.BreedImageServiceClient, args []string) {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	breed := searchCmd.String("breed", "", "Enter a breed name to search")
	subBreed := searchCmd.String("sub-breed", "", "Enter a sub-breed name to search")
	save := searchCmd.Bool("save", false, "flag to save the image to disk")
	givenPath := searchCmd.String("path", "images/", "path to save the image to")
	givenFileName := searchCmd.String("file-name", "", "file name to save the image to")

	searchCmd.Parse(args)

	log.Println("searching...")

	resp, err := c.Search(ctx, &breed_image.BreedImageSearchRequest{Breed: *breed, SubBreed: *subBreed})
	if err != nil {
		log.Printf("could not search: %v", err)
		return
	}

	if resp.ImageURL == "" || resp.Image == nil {
		log.Println("server response is not valid")
		return
	}

	if !*save {
		log.Printf("an image has been found here is the URL: \n%s\nplease add -save true flag in order to save", resp.ImageURL)
		return
	}

	log.Printf("an image has been found now saving it to disk...\n")

	fileName, err := handleFileName(*givenFileName, resp.ImageURL)
	if err != nil {
		log.Printf("could not handle file name: %v", err)
		return
	}

	fullPath, err := utils.SaveToDisk(resp.Image, fileName, *givenPath)
	if err != nil {
		log.Printf("failed to save image to disk : %v", err)
		return
	}

	str, err := filepath.Abs(fullPath)
	if err != nil {
		log.Printf("failed to get absolute path: %v", err)
		return
	}

	log.Println("image saved to disk at : ", str)

}

// helpCommand prints the help message.
func helpCommand() {
	fmt.Println("Usage: executable [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  search")
	fmt.Println("    -breed <breed> \t\t[required]")
	fmt.Println("    -sub-breed <sub-breed> \t[optional]")
	fmt.Println("    -save \t\t\t[optional]")
	fmt.Println("    -path <path> \t\t[optional]")
	fmt.Println("    -file-name <file-name> \t[optional]")
	fmt.Println("  help")
	os.Exit(1)
}

// handleFileName handles the file name as following;
// If the file name is not provided, it returns the file name from the image URL.
// If the file name is provided, it returns the file name and an error if file name is not valid.
func handleFileName(name, url string) (string, error) {
	fileNameFromURL := path.Base(url)
	if name == "" {
		return fileNameFromURL, nil
	}

	rgx := regexp.MustCompile(`[^\w-]`).MatchString
	if rgx(name) {
		return "", fmt.Errorf("file name can only word characters (letter, number, underscore[_] and dash[-])")
	}
	return fmt.Sprintf("%s.jpg", name), nil
}

// getAddr returns the address of the server.
// If the environment variable GRPC_SERVER_ADDR is set, it returns that value.
// Otherwise, it returns localhost:22626.
func getAddr() string {
	addr, ok := os.LookupEnv("CLIENT_GRPC_ADDR")
	if !ok {
		addr = "localhost:22626"
	}
	return addr
}
