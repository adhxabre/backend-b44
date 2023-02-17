# Make Hello World

### 1. Initializing project

```bash
go mod init _project_name_
```

### 2. Install gorilla/mux

```bash
go get github.com/labstack/echo/v4
```

Package `labstack/echo` implements a request router and dispatcher for matching incoming requests to their respective handler.

### 3. Create `main.go` file and write this below code to print 'Hello world'

> File : `main.go`

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.String(http.StatusOK, "Hello World")
	})

	fmt.Println("Server running on localhost:5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}
```

### 4. Running

Running Your App with this command

```
go run main.go
```
