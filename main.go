package main

import (
	"github.com/junaidshaikh-js/workout-api/internal/app"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	app.Logger.Println("Starting application...")
}
