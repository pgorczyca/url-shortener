package main

import (
	"github.com/pgorczyca/url-shortener/internal/app"
)

func main() {
	app, _ := app.NewApp()
	app.Run()

}
