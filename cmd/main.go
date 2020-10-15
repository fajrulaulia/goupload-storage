package main

import (
	app "gouploadstorage"
	"log"
)

func main() {
	log.Println("goupload-storage started on port", 8080)
	app.SetupHandler()
}
