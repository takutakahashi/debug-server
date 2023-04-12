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

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		dataMutex.RLock()
		defer dataMutex.RUnlock()

		return c.String(http.StatusOK, data)
	})

	go updateData()

	e.Start(":8080")
}
