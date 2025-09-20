package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	delstr := strings.Split(data, ",") // Разделяем строку на слайс
	if len(delstr) != 3 {
		return 0, "", 0, fmt.Errorf("Ошибка: не верный формат, ожидается 3 значение ,получено %d", len(delstr))
	}
	steps, err := strconv.Atoi(strings.TrimSpace(delstr[0])) // Парсим шаги
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка при присинге шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: кол-во шагов должно быть < 0, получено %d", steps)
	}
	activ := delstr[1]
	if activ == "" {
		return 0, "", 0, fmt.Errorf("Ошибка: неверный вид активности")
	}
	duration, err := time.ParseDuration(strings.TrimSpace(delstr[2])) //Парсим продолжительность
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка: при прасинге продолжительности %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Ошибка: продолжительность должна быть положительнаяб получено %d", duration)
	}
	return steps, activ, duration, nil
}

func distance(steps int, height float64) float64 {
	stridelength := height * stepLengthCoefficient    //Длина шага
	distance := float64(steps) * stridelength / mInKm //Дистанция в км
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	durationHours := duration.Hours()
	return distance / durationHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	if weight <= 0 {
		log.Panicln("вес должен быть положителен")
		return "", fmt.Errorf("вес должен быть положителен")
	}
	if height <= 0 {
		log.Panicln("рост должен быть положителен")
		return "", fmt.Errorf("рост должен быть положителен")
	}

	steps, activ, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// дистанция скорость
	dist := distance(steps, height)
	durationHours := duration.Hours()
	speed := meanSpeed(steps, weight, duration)

	var calorie float64

	switch activ {
	case "Бег":
		calorie, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calorie, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		log.Println("Неизвезный тип: %s", activ)
		return "", fmt.Errorf("Неизвезный тип: %s", activ)
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	// Форматируем результативную строчку
	finish := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\n Дистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activ, durationHours, dist, speed, calorie)
	return finish, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("Вес должены быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("Вес должены быть положительным")
	}
	if steps <= 0 {
		return 0, errors.New("Кол-во шагов должено быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть положительная должено быть положительным")
	}
	speed := meanSpeed(steps, height, duration)    // рассчитываем среднию скорость
	minutes := duration.Minutes()                  // переводим продолжительность в минутах
	caloris := (weight * speed * minutes) / minInH //Вычисляем калории
	return caloris, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("Вес должены быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("Вес должены быть положительным")
	}
	if steps <= 0 {
		return 0, errors.New("Кол-во шагов должено быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть положительная должено быть положительным")
	}
	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()                                               // переводим продолжительность в минутах
	caloris := (weight * speed * minutes) / minInH * walkingCaloriesCoefficient //Вычисляем калории с учетом кэф
	return caloris, nil
}
