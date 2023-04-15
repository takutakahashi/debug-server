package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	data      string
	dataMutex sync.RWMutex
	reqmap    map[string]string
	resmap    map[string]string
)

func updateData() {
	dataMutex.Lock()
	data = fmt.Sprintf("Random number: %d", rand.Intn(100))
	dataMutex.Unlock()
}

func main() {
	e := echo.New()
	reqmap = map[string]string{}
	resmap = map[string]string{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		req := c.Request()
		res := c.Response()
		if os.Getenv("CACHE") != "" {
			if req.Header.Get("if-none-match") == data {
				return c.NoContent(304)
			}
		}
		// Log request headers
		for k, v := range req.Header {
			reqmap[k] = strings.Join(v, ",")
		}
		for k, v := range res.Header() {
			resmap[k] = strings.Join(v, ",")
		}
		res.Header().Set("Etag", data)

		return c.String(http.StatusOK, data)
	})
	e.GET("/update", func(c echo.Context) error {
		updateData()
		return c.String(http.StatusOK, data)
	})
	e.GET("/headers/:t", func(c echo.Context) error {
		t := c.Param("t")
		if t == "res" {
			return c.JSON(http.StatusOK, resmap)
		} else {
			return c.JSON(http.StatusOK, reqmap)
		}
	})

	e.Start(":8080")
}
