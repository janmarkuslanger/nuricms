# nuri-cms

**nuricms is a api first content management system written in go.**

[![codecov](https://codecov.io/gh/janmarkuslanger/nuricms/graph/badge.svg?token=U51WPEFN5Y)](https://codecov.io/gh/janmarkuslanger/nuricms)
<a href="https://goreportcard.com/report/github.com/janmarkuslanger/nuricms"><img src="https://goreportcard.com/badge/github.com/janmarkuslanger/nuricms" alt="Go Report"></a>

---

<img src="demo.png" alt="Dashboard Screenshot" width="800"/>

---

## Installation & Usage

To install and use NuriCMS as a dependency in your Go project, follow the steps below.

### 1. Add NuriCMS as a Dependency

You can add NuriCMS to your Go project using `go get`. In your Go project directory, run the following command:

```bash
go get github.com/janmarkuslanger/nuricms
```

### 2. Create Your `main.go` to Start the Server

After adding NuriCMS as a dependency, you need to create a `main.go` file in your project to start the server.

#### Example `main.go`:

```go
package main

import (
	"github.com/janmarkuslanger/nuricms"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/janmarkuslanger/nuricms/pkg/config"
)

func main() {
	config := config.Config{
		Port:        "8080",
		HookPlugins: []plugin.HookPlugin{},
	}

	nuricms.Run(config)
}
```

### 3. Set JWT Secret

```bash
# Set a basic JWT secret (for development purposes)
export JWT_SECRET=anything

# Generate a secure JWT secret (recommended for production)
export JWT_SECRET=$(openssl rand -base64 32)
```

### 4. Start the Server

```bash
go run main.go
```

If the server gets started and there is no user in the system there will be an admin account added:

- E-Mail: admin@admin.com 
- Password: mysecret

The server will now run at `http://localhost:8080`. You can change the port by modifying the configuration.

---

## 🧱 How it works

At the core of **nuricms** are three key concepts:

### 1. Collections

A **collection** defines the structure of a content type – such as `blog`, `product`, or `page`. Each collection is made up of multiple fields.

### 2. Fields

A **field** describes a single property of a collection, such as `title`, `price`, or `image`. Fields have:
- A name and alias
- A type (e.g. `text`, `number`, `boolean`, `date`, `richtext`, `asset`, `collection`)
- Optional settings like default values or whether they are required

### 3. Content

Once a collection is created, you can add content entries for it. Each entry stores values for every field defined in the collection.

Values are grouped by their field alias and contain both the field value and its type.

Example: A `blog` collection with fields `title` and `body` might return the following content entry via the API:

```json
{
  "id": 29,
  "created_at": "2025-07-26T12:23:28.705057+02:00",
  "updated_at": "2025-07-26T12:23:28.705057+02:00",
  "collection": {
    "id": 12
  },
  "values": {
    "title": {
      "id": 215,
      "value": "A wonderful blog",
      "field_type": "Text"
    },
    "body": {
      "id": 216,
      "value": "<p>This is my blog body</p>",
      "field_type": "RichText"
    }
  }
}
```
---

## Plugin System

`nuricms` provides a modular plugin system. Plugins can be passed to the CMS via the `ServerConfig` at startup and allow extending various parts of the system — such as hooks, routes, or UI components.

### Usage

Plugins are passed in during server initialization:

```go
config := config.Server{
    Port: "8080",
    HookPlugins: []nuricms.HookPlugin{
        &MyCustomPlugin{},
    },
}

nuricms.Run(config)
```

### HookPlugin

A `HookPlugin` allows you to register functions for specific system events (hooks), such as "content:beforeSave". A hook plugin implements the following interface:

```go
type HookPlugin interface {
    Name() string
    Register(h *HookRegistry)
}
```

### Example

```go
package plugins

import (
	"strings"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
)

type SlugPlugin struct{}

func (p *SlugPlugin) Name() string {
	return "AutoSlug"
}

func (p *SlugPlugin) Register(h *plugin.HookRegistry) {
	h.Register("contentValue:beforeSave", func(p any) error {
		content := p.(*model.ContentValue)

		if content.Field.Alias == "slug" {
			content.Value = strings.ToLower(content.Value)
		}

		return nil
	})
}
```

Then register your plugin in your `main.go`:

```go
cfg := config.Server{
    Port: "8080",
    HookPlugins: []nuricms.HookPlugin{
        &myplugin.SlugPlugin{},
    },
}

nuricms.Run(cfg)
```

### Available Hook Events

- `contentValue:beforeSave`

---

## Architecture Overview

NuriCMS is built on a layered architecture:

```
├── cmd/              → main entry point (nuricms server)
├── internal/
│   ├── controller/   → HTTP controllers
│   ├── service/      → Business logic
│   ├── repository/   → Database access (via GORM)
│   ├── model/        → Core entities and types
│   └── handler/      → Generic HTTP handler logic (reused)
├── pkg/
│   ├── config/       → App configuration
│   └── plugin/       → Hook/plugin interface support
└── templates/        → HTML templates for admin UI
```

- Controllers handle HTTP and delegate to services (mostly via handler funcs).
- Services coordinate business logic and validate input.
- Repositories access the database.
- Plugin system allows you to register custom logic at runtime.

---

## Contributing & Development

### Code Guidelines

- All new code **must be tested**.
- Keep logic modular and covered with **unit tests**.
- Use **dependency injection** where applicable to make components testable.
- Follow standard Go formatting (`gofmt` is CI-enforced).

### Testing

To run tests and generate coverage reports:

```bash
go test ./... -coverprofile=cover.out
go tool cover -html=cover.out
```

---

For questions or contributions, feel free to open an issue or pull request.
