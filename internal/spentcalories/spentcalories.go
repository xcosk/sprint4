package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы для расчетов калорий и расстояний
const (
	lenStep                    = 0.65 // средняя длина шага в метрах
	mInKm                      = 1000 // количество метров в километре
	minInH                     = 60   // количество минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку по запятым
	delstr := strings.Split(data, ",")
	if len(delstr) != 3 {
		return 0, "", 0, fmt.Errorf("Ошибка: неверный формат, ожидается 3 значения, получено %d", len(delstr))
	}

	// Обрезаем пробелы со всех параметров
	trSteps := strings.TrimSpace(delstr[0])
	activ := strings.TrimSpace(delstr[1])
	trDuration := strings.TrimSpace(delstr[2])

	// Убираем знак + если есть перед числом шагов
	if strings.HasPrefix(trSteps, "+") {
		trSteps = trSteps[1:]
	}

	// Парсим количество шагов
	steps, err := strconv.Atoi(trSteps)
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка при парсинге шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: кол-во шагов должно быть > 0, получено %d", steps)
	}

	// Проверяем что указан тип активности
	if activ == "" {
		return 0, "", 0, fmt.Errorf("Ошибка: неверный вид активности")
	}

	// Парсим продолжительность тренировки
	duration, err := time.ParseDuration(trDuration)
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка при парсинге продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: продолжительность должна быть положительная, получено %v", duration)
	}

	return steps, activ, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага исходя из роста
	stridelength := height * stepLengthCoefficient
	// Рассчитываем общую дистанцию в километрах
	distance := float64(steps) * stridelength / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем что продолжительность положительная
	if duration <= 0 {
		return 0
	}

	// Рассчитываем дистанцию
	distance := distance(steps, height)
	// Переводим продолжительность в часы
	durationHours := duration.Hours()

	// Проверяем чтобы не было деления на ноль
	if durationHours == 0 {
		return 0
	}

	// Рассчитываем среднюю скорость
	return distance / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Проверяем корректность веса и роста
	if weight <= 0 {
		return "", fmt.Errorf("вес должен быть положителен")
	}
	if height <= 0 {
		return "", fmt.Errorf("рост должен быть положителен")
	}

	// Парсим данные тренировки
	steps, activ, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Рассчитываем дистанцию и среднюю скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calorie float64

	// В зависимости от типа активности рассчитываем калории
	switch activ {
	case "Бег":
		calorie, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calorie, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		// Если тип активности неизвестен - возвращаем ошибку
		log.Printf("Неизвестный тип тренировки: %s", activ)
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activ)
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Форматируем результат
	durationHours := duration.Hours()
	finish := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activ, durationHours, dist, speed, calorie)
	return finish, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность входных параметров
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	// Переводим продолжительность в минуты
	minutes := duration.Minutes()
	// Рассчитываем количество потраченных калорий
	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность входных параметров
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	// Переводим продолжительность в минуты
	minutes := duration.Minutes()
	// Рассчитываем количество потраченных калорий с учетом коэффициента для ходьбы
	calories := (weight * speed * minutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
