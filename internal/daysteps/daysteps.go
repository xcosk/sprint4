package daysteps

import (
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы для расчетов дистанции
const (
	stepLength = 0.65 // средняя длина одного шага в метрах
	mInKm      = 1000 // количество метров в одном километре
)

// parsePackage парсит строку с данными о шагах и продолжительности
// Формат данных: "количество_шагов,продолжительность"
// Возвращает: количество шагов, продолжительность и ошибку
func parsePackage(data string) (int, time.Duration, error) {
	delstr := strings.Split(data, ",")
	if len(delstr) != 2 {
		return 0, 0, fmt.Errorf("Ошибка: ожидалось 2 значения, получено %d", len(delstr))
	}

	// Проверяем наличие пробелов в начале или конце чисел
	trSteps := delstr[0]
	trDuration := delstr[1]

	// Если есть пробелы в начале или конце шагов - возвращаем ошибку
	if strings.HasPrefix(trSteps, " ") || strings.HasSuffix(trSteps, " ") {
		return 0, 0, fmt.Errorf("Ошибка: пробелы в количестве шагов не допускаются")
	}

	// Если есть пробелы в начале или конце продолжительности - возвращаем ошибку
	if strings.HasPrefix(trDuration, " ") || strings.HasSuffix(trDuration, " ") {
		return 0, 0, fmt.Errorf("Ошибка: пробелы в продолжительности не допускаются")
	}

	// Убираем знак + если есть
	if strings.HasPrefix(trSteps, "+") {
		trSteps = trSteps[1:]
	}

	steps, err := strconv.Atoi(trSteps)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка при парсинге шагов: %v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: кол-во шагов должно быть положительное %d", steps)
	}

	duration, err := time.ParseDuration(trDuration)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка при парсинге продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("Ошибка: продолжительность должна быть положительная, получено %s", duration)
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает данные о дневной активности и возвращает форматированную информацию
// о пройденных шагах, дистанции и потраченных калориях
func DayActionInfo(data string, weight, height float64) string {
	// Парсим данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Err: %v", err)
		return ""
	}

	// Дополнительная проверка на положительное количество шагов
	if steps <= 0 {
		return ""
	}

	// Рассчитываем пройденную дистанцию в километрах
	distance := float64(steps) * stepLength / mInKm

	// Рассчитываем потраченные калории используя функцию из пакета spentcalories
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Err: %v", err)
		return ""
	}

	// Форматируем и возвращаем результат
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories)
}
