package app

import (
	"log"
	"time"
)


func (app *Application) updateCache() {
	for {
		<-app.WebhookChan
		err := app.service.SyncData()
		if err != nil {
			log.Printf("[App] Error updating cache: %v\n", err)
		} else {
			log.Println("[App] Cache updated successfully.")
		}
		time.Sleep(1 * time.Second)
	}
}