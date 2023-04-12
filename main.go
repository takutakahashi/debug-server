package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	data      string
	dataMutex sync.RWMutex
	reqmap    map[string]string
)

func updateData() {
	for {
		dataMutex.Lock()
		data = fmt.Sprintf("Random number: %d", rand.Intn(100))
		dataMutex.Unlock()

		time.Sleep(10 * time.Minute)
	}
}

func main() {
	e := echo.New()
	reqmap = map[string]string{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		req := c.Request()

		// Log request headers
		for k, v := range req.Header {
			reqmap[k] = strings.Join(v, ",")
		}

		return c.String(http.StatusOK, data)
	})
	e.GET("/headers", func(c echo.Context) error {
		return c.JSON(http.StatusOK, reqmap)
	})

	go updateData()

	e.Start(":8080")
}
