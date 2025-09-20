package daysteps

import (
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	delstr := strings.Split(data, ",") // Разделяем строку на слайс
	if len(delstr) != 2 {
		return 0, 0, fmt.Errorf("Ошибка: ожидалось 2 значения, получено %d", len(delstr))
	}
	trSteps := strings.TrimSpace(delstr[0]) // Парсим кол-во шагов учитывая пробелы
	steps, err := strconv.Atoi(trSteps)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка при прасниге шагов: %v", err)
	}
	if steps <= 0 { // Проверяем что кол-во шагов больше 0
		fmt.Errorf("Ошибка: кол-во шагов должно быть положительное %d", steps)
	}

	duration, err := time.ParseDuration(delstr[1]) //Парсим продолжительность
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка: при прасинге продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: продолжительность должна быть положительная, получено %s", duration)
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data) //Получаем шаги и продолжительность
	if err != nil {
		log.Printf("Err: %v", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * stepLength / mInKm                                      // Вычисляем дистанцию в километрах
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration) // Вычисляем потраченные калории
	if err != nil {
		log.Printf("Err: %v", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
