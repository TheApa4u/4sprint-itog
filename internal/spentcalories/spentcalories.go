package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий ри ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Делим строку на слайс
	parts := strings.Split(data, ",")

	// Проверка на длину слайса
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: требуется 3 элемента (шаги, активность, время), получено %d", len(parts))
	}

	// Кол-во шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Невозможно распознать количество шагов: %v", err)
	}
	if steps < 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов не может быть отрицательным")
	}
	if steps == 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов не может ровняться 0")
	}

	// Преобразовывем продолжительность в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Невозможно распознать продолжительность: %v", err)
	}
	if duration < 0 {
		return 0, "", 0, fmt.Errorf("Продолжительность не может быть отрицательной")
	}
	if duration == 0 {
		return 0, "", 0, fmt.Errorf("Продолжительность не может ровняться 0")
	}

	// Для удобства
	activity := parts[1]

	// возвращаем данные и nil (если нет ошибок)
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// Общее пройденное расстояние в метрах
	totalDistanceMeters := float64(steps) * stepLength

	// Переводим метры в километры
	distanceKm := totalDistanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	//Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0.0
	}

	// Вычисляем дистанцию в км.
	distKm := distance(steps, height)

	// Продолжительность в часы
	durationHours := duration.Hours()

	// Средняя скорость
	speed := distKm / durationHours

	return speed
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка корректности входных параметров
	if steps < 0 {
		return 0, fmt.Errorf("количество шагов не может быть отрицательным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration < 0 {
		return 0, fmt.Errorf("продолжительность не может быть отрицательной")
	}

	// Расчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Формула: Каллории
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров на корректность
	if steps < 0 {
		return 0, fmt.Errorf("количество шагов не может быть отрицательным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration < 0 {
		return 0, fmt.Errorf("продолжительность не может быть отрицательной")
	}

	// Рассчитать среднюю скорость (в км/ч)
	speed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Вычислить базовое количество калорий по формуле:
	// (вес * скорость * время в минутах) / минут в часе
	calories := (weight * speed * durationMinutes) / minInH

	// Применить корректирующий коэффициент для ходьбы
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Распарсить строку с помощью parseTraining
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Проверить вес и рост
	if weight <= 0 {
		err := fmt.Errorf("вес должен быть положительным")
		log.Println(err)
		return "", err
	}
	if height <= 0 {
		err := fmt.Errorf("рост должен быть положительным")
		log.Println(err)
		return "", err
	}

	// Общие значения
	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	var calories float64

	// Обработка вида активности
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

	default:
		err := fmt.Errorf("неизвестный тип тренировки: %s", activity)
		log.Println(err)
		return "", err
	}

	result := fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f`, activity, duration.Hours(), distanceKm, speed, calories)

	return result, nil
}
