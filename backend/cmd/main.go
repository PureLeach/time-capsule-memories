// @title Time Capsule Memories API
// @version 1.0
// @description This is a sample server for Time Capsule Memories.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @license.name MIT
// @host localhost:8000
// @BasePath /

package main

import (
    "github.com/labstack/echo/v4"
    echoSwagger "github.com/swaggo/echo-swagger"
    _ "time_capsule_memories/docs"
    "time_capsule_memories/internal/routes"
)

func main() {
    e := echo.New()

    routes.RegisterRoutes(e)

    e.GET("/swagger/*", echoSwagger.WrapHandler)

    e.Logger.Fatal(e.Start(":8000"))
}
