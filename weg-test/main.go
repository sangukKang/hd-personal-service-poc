package main

import (
	"fmt"
	api "weg-test/src/api"
)


func main() {
	fmt.Println("Server start")
	api.Router()
}
