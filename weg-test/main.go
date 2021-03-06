package main

import (
	"fmt"
	"runtime"
	docs "weg-test/docs"
	api "weg-test/src/api"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Server start")
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	api.Router()
}
