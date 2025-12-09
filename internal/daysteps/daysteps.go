package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: требуется 2 значения (шаги и время), получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно распознать количество шагов: %v", err)
	}
	if steps < 0 {
		return 0, 0, fmt.Errorf("количество шагов не может быть отрицательным")
	}
	if steps == 0 {
		return 0, 0, fmt.Errorf("количество шагов не может быть нулевым")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("невозможно распознать продолжительность: %v", err)
	}
	if duration < 0 {
		return 0, 0, fmt.Errorf("продолжительность не может быть отрицательной")
	}
	if duration == 0 {
		return 0, 0, fmt.Errorf("продолжительность не может быть нулевой")
	}
	//Возвращаем значение если ошибок нет
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	//Получаем данные с помощью parsePackage()
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	//Дистанция в метрах и киллометрах
	distanceMeters := float64(steps) * stepLength // в метрах
	distanceKm := distanceMeters / mInKm          // в километрах

	//Количество сожжённых калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	//Итоговая строка:
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
