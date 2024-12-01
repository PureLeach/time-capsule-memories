package jobs

import (
	"fmt"
	"log"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/robfig/cron/v3"
)

func StartScheduler() {
	config := config.GetConfig()

	// Устанавливаем временную зону (например, UTC)
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(fmt.Sprintf("Ошибка загрузки временной зоны: %v", err))
	}

	// Создаем планировщик с заданной временной зоной
	c := cron.New(cron.WithLocation(loc))

	// Пример задачи, выполняющейся каждую минуту
	c.AddFunc(config.CronCapsuleDispatch, func() {
		JobCapsule()
	})

	// Запуск планировщика
	go func() {
		log.Println("Cron: Планировщик запущен")
		c.Start()
	}()

	// Завершение планировщика при остановке программы
	defer c.Stop()
}
