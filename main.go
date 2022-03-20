package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"path"
)

func main() {
	fmt.Println("Hello world!")
	server := echo.New()
	server.GET(path.Join("/"), Version)

	godotenv.Load()
	port := os.Getenv("PORT")

	address := fmt.Sprintf("%s:%s", "localhost", port)
	fmt.Println(address)
	server.Start(address)
}

func Version(context echo.Context) error {
	return context.JSON(http.StatusOK, map[string]interface{}{"version": 1})
}
