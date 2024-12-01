package jobs

import (
	"log"
	"time"
)

func PrintMessage() {
	log.Println("Cron: Фоновая задача запущена:", time.Now())

	log.Println("Cron: Фоновая задача завершена:", time.Now())
}
