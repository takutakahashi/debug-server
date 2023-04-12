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
	resmap    map[string]string
)

func updateData() {
	for {
		dataMutex.Lock()
		data = fmt.Sprintf("Random number: %d", rand.Intn(100))
		dataMutex.Unlock()

		time.Sleep(10 * time.Minute)
	}
}

func logHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		// Log request headers
		for k, v := range req.Header {
			reqmap[k] = strings.Join(v, ",")
		}

		err := next(c)

		// Log response headers
		for k, v := range res.Header() {
			resmap[k] = strings.Join(v, ",")
		}

		return err
	}
}

func main() {
	e := echo.New()
	reqmap = map[string]string{}
	resmap = map[string]string{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(logHeaders)

	e.GET("/", func(c echo.Context) error {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		return c.String(http.StatusOK, data)
	})
	e.GET("/headers", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]map[string]string{"req": reqmap, "res": resmap})
	})

	go updateData()

	e.Start(":8080")
}
