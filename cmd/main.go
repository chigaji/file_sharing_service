package main

import (
	"fmt"
	"log"

	"github.com/chigaji/file_sharing_service/pkg/handlers"
	md "github.com/chigaji/file_sharing_service/pkg/middleware"
	"github.com/chigaji/file_sharing_service/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("welcome")

	e := echo.New()

	//middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Routes
	e.POST("/register", handlers.RegisterHandler)
	e.POST("/login", handlers.LoginHandser)
	e.POST("/upload", handlers.UploadHandler, md.AuthMiddleware)
	e.GET("/download", handlers.DowloadHandler, md.AuthMiddleware)

	//initialize db
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Error Initializing datatabase: %v", err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
