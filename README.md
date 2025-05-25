# nuri-cms

**nuricms is a api first content management system written in go.**

---

![Coverage Badge](coverage_badge.svg)

---

## Installation

To install and use NuriCMS as a dependency in your Go project, follow the steps below.

### 1. Add NuriCMS as a Dependency

You can add NuriCMS to your Go project using `go get`. In your Go project directory, run the following command:

```bash
go get github.com/janmarkuslanger/nuricms
```

This will add NuriCMS as a dependency to your project and fetch the latest version.

### 2. Create Your `main.go` to Start the Server

After adding NuriCMS as a dependency, you need to create a `main.go` file in your project to start the server.

#### Example `main.go`:

```go
package main

import (
	"log"
	"os"
	"github.com/janmarkuslanger/nuricms"
)

func main() {
	config := &nuricms.ServerConfig{
		Port: os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	nuricms.StartServer(config)
}
```

Then run the following command to start the server:

```bash
go run main.go
```

The server will now run at `http://localhost:8080`. You can change the port by modifying the configuration.

## Docker 

You can find every dockerfile in folder `/docker`. 

### Build the container

If you want to build one of the dockerfiles you need to enter: 

`docker build -t nuricms -f path/to/dockerfile .`

For example: 
`docker build -t nuricms -f docker/nuricms-sqlite/Dockerfile .`

### Run the container

`docker run -p 8080:8080 -it -e JWT_SECRET=my-verysuper-secret-secret-32byteslong nuricms`