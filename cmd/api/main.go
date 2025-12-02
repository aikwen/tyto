package main

import (
	"tyto/internal/app"
)

func main() {
	application := app.NewApplication()
	application.RegisterRoutes()
	application.Run(":9001")
}

