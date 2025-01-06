package main

import (
	"abacus_engine/internal/controller"
)

func main() {
	rows := 5000
	appController := controller.NewAppController(rows)

	if err := appController.Run(); err != nil {
		panic(err)
	}
}
