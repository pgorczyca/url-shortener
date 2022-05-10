package main

import (
	"github.com/pgorczyca/url-shortener/internal/app"
)

func main() {
	app, _ := app.NewApp()
	app.Run()
	// url1 := url{
	// 	Long:      "https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/",
	// 	Short:     "ASt",
	// 	ExpiredAt: time.Now().Add(time.Hour * 6),
	// 	CreatedAt: time.Now(),
	// }
	// repo.add(context.TODO(), url1)
}
