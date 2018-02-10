package main

import (
	//"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"io/ioutil"
	//"log"
	"net/http"
)


func mainAdmin(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
			"id": "1",
			"name": "Grow with Google",
		})


}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()
	g := e.Group("/")

	// this logs the server interaction
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "google" && password == "group" {
			return true, nil

		}

		return false, nil

	}))

	g.POST("basic", mainAdmin)

	e.Start(":3000")
}
