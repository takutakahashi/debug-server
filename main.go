package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	data      string
	dataMutex sync.RWMutex
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
		var reqHeaders []string
		for k, v := range req.Header {
			reqHeaders = append(reqHeaders, fmt.Sprintf("%s: %s", k, strings.Join(v, ",")))
		}
		c.Logger().Infof("Request headers:\n%s", strings.Join(reqHeaders, "\n"))

		err := next(c)

		// Log response headers
		var resHeaders []string
		for k, v := range res.Header() {
			resHeaders = append(resHeaders, fmt.Sprintf("%s: %s", k, strings.Join(v, ",")))
		}
		c.Logger().Infof("Response headers:\n%s", strings.Join(resHeaders, "\n"))

		return err
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
  e.Use(logHeaders)

	e.GET("/", func(c echo.Context) error {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		return c.String(http.StatusOK, data)
	})

	go updateData()

	e.Start(":8080")
}
