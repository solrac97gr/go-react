package main

import (
	"io/fs"
	"log"

	app "github.com/solrac97gr/go-react"
	"github.com/solrac97gr/go-react/pkg/server"
)

var uiFS fs.FS

func init() {
	var err error
	uiFS, err = fs.Sub(app.UI, "_ui/build")
	if err != nil {
		log.Fatal("failed to get ui fs", err)
	}
}

func main() {
	if err := server.InitServer(uiFS, ":8080"); err != nil {
		log.Fatal(err)
	}
}
