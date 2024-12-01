package jobs

import (
	"log"
	"sync"
	"time"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/services"
)

func JobCapsule() {
	log.Println("Cron: Фоновая задача запущена:", time.Now())

	// Вызов метода для получения капсул за текущую дату
	capsules, err := repository.GetCapsulesByToday()
	if err != nil {
		log.Printf("Ошибка при получении данных: %v", err)
		return
	}

	// Если список пуст, ничего не делаем
	if len(capsules) == 0 {
		log.Println("Cron: Нет данных для обработки")
	} else {
		// Создаем WaitGroup для ожидания завершения всех горутин
		var wg sync.WaitGroup

		// Если данные есть, выводим id каждой записи и запускаем горутину для обработки каждой капсулы
		for _, capsule := range capsules {
			wg.Add(1) // Увеличиваем счетчик горутин
			go func(capsule models.CapsuleResponse) {
				defer wg.Done()                   // Уменьшаем счетчик при завершении горутины
				services.ProcessCapsule(&capsule) // Pass pointer to the function
			}(*capsule) // Passing the value of capsule directly
		}

		// Ожидаем завершения всех горутин
		wg.Wait()
	}

	log.Println("Cron1: Фоновая задача завершена:", time.Now())
}
