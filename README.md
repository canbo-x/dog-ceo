# dog-ceo
This is a project to demonstrate a basic gRPC server in Go.
There are server and client sides. Both the client application and backend service are written in Go.

The server side is responsible for communicating with the client-side and fetching data from a public API. It fetches the data from the API over HTTP and serves this data to the client over gRPC and gets the required parameters from the client.

You can find this [dog ceo public api](https://dog.ceo/dog-api/) in the link. We used only the `dog.ceo/api/breeds/image/random` endpoint to fetch image URLs by breed and sub-breed.

The responsibility of the client application is to take breed and sub-breed from a user and make a request to the server. The expected response from the server is the image of the given input. The client is responsible for saving the image to the filesystem and letting the user know where it is saved. The user is informed if any error occurs.


# Dependencies
It is developed mostly with go standard packages. (*.pb.go and test files are excluded)
### Used packages
 ``` 
  "context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "time"
  "flag"
  "log"
  "net"
  "net/http"
  "regexp"
  "os"
  "path"
  "path/filepath"
  "time"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
```

# Installation and Usage
You must have golang installed on your computer to run the project.

Clone the project and compile.
Please visit [official doc](https://go.dev/doc/tutorial/compile-install) for more information.
```shell
git@github.com:canbo-x/dog-ceo.git
```

You need to run the server first. Navigate to the project directory in your terminal and run
```shell
./grpc_server
```

You can add a port flag to run the server in the desired port. The default port is `22626`.

```shell
./grpc_server -port 12345
```

You can set a log-level flag. The default log-level is `info`.
```shell
./grpc_server -log-level debug
```

Available options are `debug`, `info`, `warn`, `error`, `fatal`, `panic`, `trace`.

---

After the server is running you can run the client.
Before you run please see the help command.

```shell
./grpc_client -help
```

```
Usage: executable [command] [flags]
Commands:
search
  -breed <breed> [required]
  -sub-breed <sub-breed> [optional]
  -save [optional]
  -path <path> [optional]
  -file-name <file-name> [optional]
-help
```

Examples:
```shell
./grpc_client search -breed husky

./grpc_client search -breed wolfhound -sub-breed irish -save

./grpc_client search -breed wolfhound -sub-breed irish -save -path images/ -file-name lovelyDog
```

`-save` flag is required to save the image.

The default address is `localhost:22626`. You can set the environmental variable to change.
```shell
export CLIENT_GRPC_ADDR="localhost:22626" && echo $CLIENT_GRPC_ADDR
```

# Testing
```shell
go test ./...
```

```shell
ok  	github.com/canbo-x/dog-ceo/breed_image_service	3.350s
ok  	github.com/canbo-x/dog-ceo/cmd/grpc_client	3.266s
ok  	github.com/canbo-x/dog-ceo/cmd/grpc_server	9.495s
ok  	github.com/canbo-x/dog-ceo/data_service	8.156s
ok  	github.com/canbo-x/dog-ceo/dummy_rate_limiter	10.402s
ok  	github.com/canbo-x/dog-ceo/utils	0.207s
```

# Personal Thoughts and Notes

- We could also use `dog.ceo/api/breeds/list/all` to fetch available breeds and sub-breeds and cache them with the [go-cache](https://github.com/patrickmn/go-cache) package. It would make a huge difference regarding the cost and time. Because we don't have to go to the server every single time to know if the given breed or sub-breed exists. We could bind the breeds to a map so the complexity of searching would be o1.

- We could also cache the images on the server to reduce the network cost, but it is much better to use CDN because it might rapidly consume the server's storage.  Therefore increasing the cost of maintenance and storage.

- Caching the Image URLs would help to reduce the cost and response time. But in this case, we would have to handle the random mechanism and fetch all the image URLs. Honestly, I’m not so sure about this trade-off.

- Creating an HTTP Client for each request is okay for this demonstration but in a real-world example, we should reuse the client.

- We could add retry logic for the HTTP Client.

- We could inform the backend whenever a photo is successfully saved to the client machine and log it. we could also log the saving errors to have more observability. For instance, we could log errors with some info like operating system, available disk space, etc. Imagine that there is an issue with Windows OS, so we could see that there are lots of errors from a specific OS, and check my code for it.

- We could add a health check endpoint to the server.

- Please note that the dummy rate limiting service is just a demonstration. It is not intended to use in any production environment.




