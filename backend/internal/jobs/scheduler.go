package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	// Устанавливаем временную зону (например, UTC)
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(fmt.Sprintf("Ошибка загрузки временной зоны: %v", err))
	}

	// Создаем планировщик с заданной временной зоной
	c := cron.New(cron.WithLocation(loc))

	// Добавляем задачу, которая выполняется каждый день в 3:00 (в UTC)
	c.AddFunc("0 3 * * *", func() {
		PrintMessage()
	})

	// Пример задачи, выполняющейся каждую минуту
	c.AddFunc("*/1 * * * *", func() {
		PrintMessage()
	})

	// Запуск планировщика
	go func() {
		log.Println("Cron: Планировщик запущен")
		c.Start()
	}()

	// Завершение планировщика при остановке программы
	defer c.Stop()
}
